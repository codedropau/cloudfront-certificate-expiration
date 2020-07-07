package message

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/gosuri/uitable"

	"github.com/codedropau/cloudfront-certificate-expiration/internal/discovery"
)

// Send an "about to expire" message to SNS topic.
func Send(client snsiface.SNSAPI, topic string, certificates []discovery.Certificate) (string, error) {
	table := uitable.New()
	table.MaxColWidth = 80

	table.AddRow("CLOUDFRONT", "NAME", "TYPE", "EXPIRATION")

	for _, certificate := range certificates {
		table.AddRow(certificate.CloudFront, certificate.Name, certificate.Type, certificate.Expiration)
	}

	input := &sns.PublishInput{
		Subject:  aws.String("CloudFront certificates are about to expire"),
		Message:  aws.String(table.String()),
		TopicArn: aws.String(topic),
	}

	resp, err := client.Publish(input)
	if err != nil {
		return "", err
	}

	return *resp.MessageId, nil
}
