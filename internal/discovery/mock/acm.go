package mock

import (
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/acm/acmiface"
)

// ACM client for mock interactions.
type ACM struct {
	acmiface.ACMAPI
	Output *acm.DescribeCertificateOutput
}

// DescribeCertificate for mock interactions.
func (a *ACM) DescribeCertificate(*acm.DescribeCertificateInput) (*acm.DescribeCertificateOutput, error) {
	return a.Output, nil
}
