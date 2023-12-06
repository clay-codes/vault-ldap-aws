package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	// Initialize a session in us-west-2 that the SDK will use to load credentials
	// from the shared credentials file ~/.aws/credentials.
	output, err := Auth()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Output:", output)
	
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create an EC2 service client.
	svc := ec2.New(sess)

	// Call to get detailed information on each instance
	result, err := svc.DescribeInstances(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through the instances, and print their IDs
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Printf("Instance ID: %s\n", *instance.InstanceId)
		}
	}
}

// run the main function
