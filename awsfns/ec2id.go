package awsfns

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func GetEC2ID(sess *session.Session) {
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
