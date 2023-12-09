package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/clay-codes/aws-ldap/cloud"
)

func main() {
	// Initialize a session in us-west-2 that the SDK will use to load credentials
	// from the shared credentials file ~/.aws/credentials.
	if _, err := cloud.Auth(); err != nil {
		fmt.Println("Error authing: ")
		log.Fatal(err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		log.Fatal(err)
	}

	imgID, err := cloud.GetImgID(sess)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(imgID)

	str, err := cloud.CreateKP(sess)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	err = CreateSG(sess, groupName, vpcID)
	if err != nil {
		panic(err)
	}

}

// run the main function
