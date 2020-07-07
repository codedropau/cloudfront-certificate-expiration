package mock

import (
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
)

// CloudFront client for mock interactions.
type CloudFront struct {
	cloudfrontiface.CloudFrontAPI
	Output *cloudfront.ListDistributionsOutput
}

// ListDistributions for mock interactions.
func (c *CloudFront) ListDistributions(*cloudfront.ListDistributionsInput) (*cloudfront.ListDistributionsOutput, error) {
	return c.Output, nil
}
