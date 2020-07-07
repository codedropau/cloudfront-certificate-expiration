package certutils

import (
	"time"

	"github.com/codedropau/cloudfront-certificate-expiration/internal/discovery"
)

// Filter recent certificates.
func Filter(certificates []discovery.Certificate, line time.Time) []discovery.Certificate {
	var filtered []discovery.Certificate

	for _, certificate := range certificates {
		if certificate.Expiration.Before(line) {
			filtered = append(filtered, certificate)
		}
	}

	return filtered
}
