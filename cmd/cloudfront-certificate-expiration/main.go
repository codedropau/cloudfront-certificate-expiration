package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sns"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/codedropau/cloudfront-certificate-expiration/internal/discovery"
	"github.com/codedropau/cloudfront-certificate-expiration/internal/message"
)

var cliTopicARN = kingpin.Flag("topic", "AWS SNS topic ARN to post messages").Envar("CLOUDFRONT_CERTIFICATE_EXPIRATION_TOPIC_ARN").String()

func main() {
	kingpin.Parse()

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	client, err := discovery.New(cloudfront.New(sess), acm.New(sess), iam.New(sess))
	if err != nil {
		panic(err)
	}

	certificates, err := client.GetCertificates()
	if err != nil {
		panic(err)
	}

	id, err := message.Send(sns.New(sess), *cliTopicARN, certificates)
	if err != nil {
		panic(err)
	}

	fmt.Println("Message sent with ID:", id)
}
