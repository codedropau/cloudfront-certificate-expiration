package discovery

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/assert"

	"github.com/codedropau/cloudfront-certificate-expiration/internal/discovery/mock"
)

func TestGetCertificates(t *testing.T) {
	now := time.Now()

	cfclient := &mock.CloudFront{
		Output: &cloudfront.ListDistributionsOutput{
			DistributionList: &cloudfront.DistributionList{
				IsTruncated: aws.Bool(false),
				Items: []*cloudfront.DistributionSummary{
					{
						Id: aws.String("cloudfront-1"),
						ViewerCertificate: &cloudfront.ViewerCertificate{
							ACMCertificateArn: aws.String("acm-certificate-1"),
						},
					},
					{
						Id: aws.String("cloudfront-2"),
						ViewerCertificate: &cloudfront.ViewerCertificate{
							IAMCertificateId: aws.String("iam-certificate-2"),
						},
					},
				},
			},
		},
	}

	iamclient := &mock.IAM{
		Output: &iam.ListServerCertificatesOutput{
			IsTruncated: aws.Bool(false),
			ServerCertificateMetadataList: []*iam.ServerCertificateMetadata{
				{
					ServerCertificateName: aws.String("example_com"),
					ServerCertificateId:   aws.String("iam-certificate-2"),
					Expiration:            aws.Time(now),
				},
			},
		},
	}

	acmclient := &mock.ACM{
		Output: &acm.DescribeCertificateOutput{
			Certificate: &acm.CertificateDetail{
				DomainName:     aws.String("example.com"),
				CertificateArn: aws.String("acm-certificate-1"),
				NotAfter:       aws.Time(now),
			},
		},
	}

	client, err := New(cfclient, acmclient, iamclient)
	assert.Nil(t, err)

	certificates, err := client.GetCertificates()
	assert.Nil(t, err)

	want := []Certificate{
		{
			Name:       "example.com",
			CloudFront: "cloudfront-1",
			Type:       CertificateTypeACM,
			Expiration: now,
		},
		{
			Name:       "example_com",
			CloudFront: "cloudfront-2",
			Type:       CertificateTypeIAM,
			Expiration: now,
		},
	}

	assert.Equal(t, want, certificates)
}
