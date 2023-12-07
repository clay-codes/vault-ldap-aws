package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/clay-codes/aws-ldap/awsfns"
)

func main() {
	// Initialize a session in us-west-2 that the SDK will use to load credentials
	// from the shared credentials file ~/.aws/credentials.
	if _, err := awsfns.Auth(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		log.Fatal(err)
	}
	awsfns.GetEC2ID(sess)

}

// run the main function
