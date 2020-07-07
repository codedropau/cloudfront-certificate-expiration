package message

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/codedropau/cloudfront-certificate-expiration/internal/discovery"
	"github.com/codedropau/cloudfront-certificate-expiration/internal/message/mock"
)

func TestSend(t *testing.T) {
	client := &mock.SNS{}

	certificates := []discovery.Certificate{
		{
			Name:       "certificate-1",
			CloudFront: "cloudfront-1",
			Type:       discovery.CertificateTypeACM,
			Expiration: time.Now().Add(5 * 24 * time.Hour),
		},
		{
			Name:       "certificate-2",
			CloudFront: "cloudfront-2",
			Type:       discovery.CertificateTypeACM,
			Expiration: time.Now().Add(15 * 24 * time.Hour),
		},
	}

	_, err := Send(client, "xxxxxxxxxxxxx", certificates)
	assert.Nil(t, err)

	assert.Contains(t, *client.Input.Message, "certificate-1")
	assert.Contains(t, *client.Input.Message, "certificate-2")
}
