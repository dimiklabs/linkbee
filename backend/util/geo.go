package util

import (
	"encoding/json"
	"net"
	"net/http"
	"time"
)

var geoHTTPClient = &http.Client{Timeout: 2 * time.Second}

type geoAPIResponse struct {
	Status      string `json:"status"`
	CountryCode string `json:"countryCode"`
	City        string `json:"city"`
}

// LookupGeo returns the ISO 3166-1 alpha-2 country code and city for the given IP.
// Returns empty strings for private/loopback IPs or on any lookup failure (graceful degradation).
func LookupGeo(ip string) (countryCode, city string) {
	if isPrivateOrLoopback(ip) {
		return "", ""
	}

	resp, err := geoHTTPClient.Get("http://ip-api.com/json/" + ip + "?fields=status,countryCode,city")
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()

	var geo geoAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil || geo.Status != "success" {
		return "", ""
	}

	return geo.CountryCode, geo.City
}

func isPrivateOrLoopback(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return true
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast()
}
