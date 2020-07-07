package discovery

import (
	"time"

	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/acm/acmiface"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

// Client for discovering certificates.
type Client struct {
	CloudFront cloudfrontiface.CloudFrontAPI
	ACM        acmiface.ACMAPI
	IAM        iamiface.IAMAPI
}

// CertificateType identifies which managed service the certificate is stored.
type CertificateType string

const (
	// CertificateTypeACM is used to identify certificates provisioned with ACM.
	CertificateTypeACM = "ACM"
	// CertificateTypeIAM is used to identify certificates stored in IAM.
	CertificateTypeIAM = "IAM"
)

// Certificate information.
type Certificate struct {
	Name       string
	CloudFront string
	Type       CertificateType
	Expiration time.Time
}

// New client for discovering certificates.
func New(cloudfront cloudfrontiface.CloudFrontAPI, acm acmiface.ACMAPI, iam iamiface.IAMAPI) (*Client, error) {
	return &Client{cloudfront, acm, iam}, nil
}

// GetCertificates which are associated with CloudFront distributions.
func (c *Client) GetCertificates() ([]Certificate, error) {
	var certificates []Certificate

	distributions, err := getDistributions(c.CloudFront)
	if err != nil {
		return certificates, err
	}

	servers, err := getServerCertificates(c.IAM)
	if err != nil {
		return certificates, err
	}

	for _, distribution := range distributions {
		if distribution.ViewerCertificate == nil {
			continue
		}

		if distribution.ViewerCertificate.ACMCertificateArn != nil {
			resp, err := c.ACM.DescribeCertificate(&acm.DescribeCertificateInput{
				CertificateArn: distribution.ViewerCertificate.ACMCertificateArn,
			})
			if err != nil {
				return nil, err
			}

			certificates = append(certificates, Certificate{
				Name:       *resp.Certificate.DomainName,
				CloudFront: *distribution.Id,
				Type:       CertificateTypeACM,
				Expiration: *resp.Certificate.NotAfter,
			})

			continue
		}

		if distribution.ViewerCertificate.IAMCertificateId != nil {
			for _, server := range servers {
				if *server.ServerCertificateId != *distribution.ViewerCertificate.IAMCertificateId {
					continue
				}

				certificates = append(certificates, Certificate{
					Name:       *server.ServerCertificateName,
					CloudFront: *distribution.Id,
					Type:       CertificateTypeIAM,
					Expiration: *server.Expiration,
				})
			}
		}
	}

	return certificates, nil
}

// Helper function to get all CloudFront distributions and handle pagination.
func getDistributions(client cloudfrontiface.CloudFrontAPI) ([]*cloudfront.DistributionSummary, error) {
	var list []*cloudfront.DistributionSummary

	params := &cloudfront.ListDistributionsInput{}

	for {
		resp, err := client.ListDistributions(params)
		if err != nil {
			return nil, err
		}

		list = append(list, resp.DistributionList.Items...)

		if !*resp.DistributionList.IsTruncated {
			break
		}

		params.Marker = resp.DistributionList.Marker
	}

	return list, nil
}

// Helper function to get server certificates and handle pagination.
func getServerCertificates(client iamiface.IAMAPI) ([]*iam.ServerCertificateMetadata, error) {
	var list []*iam.ServerCertificateMetadata

	params := &iam.ListServerCertificatesInput{}

	for {
		resp, err := client.ListServerCertificates(params)
		if err != nil {
			return nil, err
		}

		list = append(list, resp.ServerCertificateMetadataList...)

		if !*resp.IsTruncated {
			break
		}

		params.Marker = resp.Marker
	}

	return list, nil
}
