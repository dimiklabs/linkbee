package util

import (
	"crypto/sha256"
	"fmt"
	"net"
	"strings"
)

const ipHashSalt = "shortlink-ip-salt-v1"

// HashIP creates a salted SHA-256 hash of an IP address for GDPR-compliant storage.
func HashIP(ip string) string {
	// Normalize IP
	ip = strings.TrimSpace(ip)
	if parsed := net.ParseIP(ip); parsed != nil {
		ip = parsed.String()
	}

	h := sha256.New()
	h.Write([]byte(ipHashSalt + ip))
	return fmt.Sprintf("%x", h.Sum(nil))
}
