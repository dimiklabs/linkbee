package geo

import (
	"net"
	"strings"

	"go.uber.org/zap"

	"github.com/oschwald/geoip2-golang"
	"github.com/shafikshaon/linkbee/logger"
)

// GeoServiceI resolves a visitor's country code from request context.
type GeoServiceI interface {
	// GetCountryCode returns an ISO 3166-1 alpha-2 code (e.g. "US", "GB").
	// It checks CDN headers first, then the MaxMind DB (if configured).
	// Returns "" if the country cannot be determined.
	GetCountryCode(ip string, headers map[string]string) string
}

type geoService struct {
	reader *geoip2.Reader // may be nil if no DB is configured
}

// NewGeoService creates a GeoService. If geoDBPath is empty or the file cannot
// be opened the service still works — it falls back to CDN / proxy headers.
func NewGeoService(geoDBPath string) GeoServiceI {
	svc := &geoService{}
	if geoDBPath != "" {
		reader, err := geoip2.Open(geoDBPath)
		if err != nil {
			logger.Warn("GeoIP DB not loaded — header-based fallback only",
				zap.String("path", geoDBPath), zap.Error(err))
		} else {
			logger.Info("GeoIP DB loaded", zap.String("path", geoDBPath))
			svc.reader = reader
		}
	}
	return svc
}

// GetCountryCode resolves the country code for an incoming request.
// Priority:
//  1. CDN / reverse-proxy headers (CF-IPCountry, X-GeoIP-Country, X-Country-Code)
//  2. MaxMind GeoLite2 DB lookup (if configured)
func (s *geoService) GetCountryCode(ip string, headers map[string]string) string {
	// 1. Check CDN / Nginx GeoIP headers
	for _, header := range []string{"CF-IPCountry", "X-GeoIP-Country", "X-Country-Code"} {
		if v, ok := headers[header]; ok {
			v = strings.TrimSpace(strings.ToUpper(v))
			// "XX" is Cloudflare's sentinel for unknown country
			if v != "" && v != "XX" && len(v) == 2 {
				return v
			}
		}
	}

	// 2. MaxMind DB lookup
	if s.reader != nil && ip != "" {
		parsedIP := net.ParseIP(ip)
		if parsedIP != nil {
			record, err := s.reader.Country(parsedIP)
			if err == nil && record.Country.IsoCode != "" {
				return strings.ToUpper(record.Country.IsoCode)
			}
		}
	}

	return ""
}
