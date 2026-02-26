package util

import (
	"strings"
)

type DeviceInfo struct {
	DeviceType string
	DeviceName string
	Browser    string
	OS         string
}

const (
	DeviceTypeMobile  = "mobile"
	DeviceTypeTablet  = "tablet"
	DeviceTypeDesktop = "desktop"
	DeviceTypeUnknown = "unknown"
)

func ParseUserAgent(userAgent string) *DeviceInfo {
	info := &DeviceInfo{
		DeviceType: DeviceTypeUnknown,
		DeviceName: "Unknown Device",
		Browser:    "Unknown",
		OS:         "Unknown",
	}

	if userAgent == "" {
		return info
	}

	ua := strings.ToLower(userAgent)
	info.OS = detectOS(ua, userAgent)
	info.Browser = detectBrowser(ua, userAgent)
	info.DeviceType = detectDeviceType(ua)
	info.DeviceName = generateDeviceName(info)

	return info
}

func detectOS(ua, originalUA string) string {
	switch {
	case strings.Contains(ua, "windows nt 10"):
		return "Windows 10/11"
	case strings.Contains(ua, "windows nt 6.3"):
		return "Windows 8.1"
	case strings.Contains(ua, "windows nt 6.1"):
		return "Windows 7"
	case strings.Contains(ua, "windows"):
		return "Windows"
	case strings.Contains(ua, "mac os x"):
		if strings.Contains(ua, "iphone") {
			return "iOS"
		}
		if strings.Contains(ua, "ipad") {
			return "iPadOS"
		}
		return "macOS"
	case strings.Contains(ua, "android"):
		return "Android"
	case strings.Contains(ua, "linux"):
		return "Linux"
	case strings.Contains(ua, "cros"):
		return "Chrome OS"
	default:
		return "Unknown"
	}
}

func detectBrowser(ua, originalUA string) string {
	switch {
	case strings.Contains(ua, "edg/"):
		return "Edge"
	case strings.Contains(ua, "opr/") || strings.Contains(ua, "opera"):
		return "Opera"
	case strings.Contains(ua, "chrome") && !strings.Contains(ua, "chromium"):
		return "Chrome"
	case strings.Contains(ua, "firefox"):
		return "Firefox"
	case strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome"):
		return "Safari"
	case strings.Contains(ua, "msie") || strings.Contains(ua, "trident"):
		return "Internet Explorer"
	default:
		return "Unknown"
	}
}

func detectDeviceType(ua string) string {
	switch {
	case strings.Contains(ua, "mobile") || strings.Contains(ua, "iphone"):
		return DeviceTypeMobile
	case strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad"):
		return DeviceTypeTablet
	case strings.Contains(ua, "windows") || strings.Contains(ua, "macintosh") || strings.Contains(ua, "linux"):
		return DeviceTypeDesktop
	default:
		return DeviceTypeUnknown
	}
}

func generateDeviceName(info *DeviceInfo) string {
	if info.Browser != "Unknown" && info.OS != "Unknown" {
		return info.Browser + " on " + info.OS
	}
	if info.Browser != "Unknown" {
		return info.Browser
	}
	if info.OS != "Unknown" {
		return info.OS + " Device"
	}
	return "Unknown Device"
}
