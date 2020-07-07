package mock

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

// SNS for mock interactions.
type SNS struct {
	snsiface.SNSAPI
	Input *sns.PublishInput
}

// Publish for mock interactions.
func (s *SNS) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	s.Input = input

	return &sns.PublishOutput{
		MessageId: aws.String("some-message-id"),
	}, nil
}
