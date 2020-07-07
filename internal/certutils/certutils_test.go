package certutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/codedropau/cloudfront-certificate-expiration/internal/discovery"
)

func TestFilter(t *testing.T) {
	now := time.Now()

	certificates := []discovery.Certificate{
		{
			Name:       "certificate-1",
			CloudFront: "cloudfront-1",
			Type:       discovery.CertificateTypeACM,
			Expiration: now.Add(5 * 24 * time.Hour),
		},
		{
			Name:       "certificate-2",
			CloudFront: "cloudfront-2",
			Type:       discovery.CertificateTypeACM,
			Expiration: now.Add(15 * 24 * time.Hour),
		},
		{
			Name:       "certificate-2",
			CloudFront: "cloudfront-2",
			Type:       discovery.CertificateTypeACM,
			Expiration: now.Add(30 * 24 * time.Hour),
		},
		{
			Name:       "certificate-3",
			CloudFront: "cloudfront-3",
			Type:       discovery.CertificateTypeACM,
			Expiration: now.Add(60 * 24 * time.Hour),
		},
		{
			Name:       "certificate-4",
			CloudFront: "cloudfront-4",
			Type:       discovery.CertificateTypeACM,
			Expiration: now.Add(90 * 24 * time.Hour),
		},
	}

	want := []discovery.Certificate{
		{
			Name:       "certificate-1",
			CloudFront: "cloudfront-1",
			Type:       discovery.CertificateTypeACM,
			Expiration: now.Add(5 * 24 * time.Hour),
		},
		{
			Name:       "certificate-2",
			CloudFront: "cloudfront-2",
			Type:       discovery.CertificateTypeACM,
			Expiration: now.Add(15 * 24 * time.Hour),
		},
	}

	assert.Equal(t, want, Filter(certificates, now.Add(30*24*time.Hour)))
}
