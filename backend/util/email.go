package util

import (
	"strings"
)

// gmailDomains are providers that ignore dots in the local part of an email address.
var gmailDomains = map[string]bool{
	"gmail.com":      true,
	"googlemail.com": true,
}

// NormalizeEmail returns a canonical form of the email address so that
// addresses which are equivalent from the provider's perspective are stored
// and compared identically.
//
// Rules applied:
//  1. Lowercase the entire address.
//  2. Trim leading/trailing whitespace.
//  3. For Gmail / Googlemail: remove dots from the local part
//     (Gmail ignores dots, so foo.bar@gmail.com == foobar@gmail.com).
func NormalizeEmail(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))

	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return email
	}

	local, domain := parts[0], parts[1]

	if gmailDomains[domain] {
		local = strings.ReplaceAll(local, ".", "")
	}

	return local + "@" + domain
}
