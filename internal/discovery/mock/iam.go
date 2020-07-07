package mock

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

// IAM client for mock interactions.
type IAM struct {
	iamiface.IAMAPI
	Output *iam.ListServerCertificatesOutput
}

// ListServerCertificates for mock interactions.
func (i *IAM) ListServerCertificates(input *iam.ListServerCertificatesInput) (*iam.ListServerCertificatesOutput, error) {
	return i.Output, nil
}
