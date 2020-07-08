package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sns"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/codedropau/cloudfront-certificate-expiration/internal/certutils"
	"github.com/codedropau/cloudfront-certificate-expiration/internal/discovery"
	"github.com/codedropau/cloudfront-certificate-expiration/internal/message"
)

var (
	cliTopicARN = kingpin.Flag("topic", "AWS SNS topic ARN to post messages").Envar("CLOUDFRONT_CERTIFICATE_EXPIRATION_TOPIC_ARN").String()
	cliAge      = kingpin.Flag("age", "How long before development teams should be notified").Default("720h").Envar("CLOUDFRONT_CERTIFICATE_EXPIRATION_AGE").Duration()
)

func main() {
	kingpin.Parse()
	lambda.Start(HandleRequest)
}

// HandleRequest contains the code which will be executed.
func HandleRequest() error {
	// Session for interacting with global sessions (CloudFront/ACM/IAM).
	global, err := session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	})
	if err != nil {
		return err
	}

	// Session for interacting with local services (SNS).
	local, err := session.NewSession()
	if err != nil {
		return err
	}

	client, err := discovery.New(cloudfront.New(global), acm.New(global), iam.New(global))
	if err != nil {
		return err
	}

	certificates, err := client.GetCertificates()
	if err != nil {
		return err
	}

	filtered := certutils.Filter(certificates, time.Now().Add(*cliAge))

	if len(filtered) == 0 {
		return nil
	}

	fmt.Println("Sending message to SNS topic arn:", *cliTopicARN)

	id, err := message.Send(sns.New(local), *cliTopicARN, filtered)
	if err != nil {
		return err
	}

	fmt.Println("Message sent with ID:", id)

	return nil
}
