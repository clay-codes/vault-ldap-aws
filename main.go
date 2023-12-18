package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/clay-codes/aws-ldap/cloud"
)

// var EC2ID = ""
func bootStrap() {
	str, err := cloud.CreateKP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	_, err = cloud.CreateSG([]int64{22, 8200, 8201})
	if err != nil {
		log.Fatal(err)
	}

	err = cloud.CreateInstProf()
	if err != nil {
		log.Fatal(err)
	}
	// wait for instance profile to be created
	time.Sleep(10 * time.Second)

	ec2ID, err := cloud.BuildEC2()
	// EC2ID = ec2ID
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("EC2ID: %s", ec2ID)
}

func CleanupCloud() {
	if err := os.Remove("key.pem"); err != nil {
		fmt.Printf("error deleting key.pem file: %v", err)
	}
	if err := cloud.TerminateEC2Instance(); err != nil {
		fmt.Printf("error terminating EC2 instance: %v", err)
	}
	if err := cloud.DeleteKeyPair(); err != nil {
		fmt.Printf("error deleting key pair: %v", err)
	}
	if err := cloud.DetachRoleFromInstanceProfile(); err != nil {
		fmt.Printf("error detaching role from instance profile: %v", err)
	}
	if err := cloud.DeleteInstanceProfile(); err != nil {
		fmt.Printf("error deleting instance profile: %v", err)
	}
	if err := cloud.DeleteRole(); err != nil {
		fmt.Printf("error deleting role: %v", err)
	}
	if err := cloud.DeleteSecurityGroup(); err != nil {
		fmt.Printf("error deleting security group: %v", err)
	}
}

// run the main function
func main() {
	// Initialize a session in us-west-2 that the SDK will use to load credentials
	// from the shared credentials file ~/.aws/credentials.

	//set flag to run cleanup
	cloud.CheckAuth()

	// creating a session
	if err := cloud.CreateSession("us-west-2"); err != nil {
		log.Fatal(err)
	}

	// getting session
	// sess := cloud.GetSession().GetAWSSession()

	// creating needed services from session
	if err := cloud.GetSession().CreateServices(); err != nil {
		log.Fatal(err)
	}

	cleanupFlag := flag.Bool("cleanup", false, "Set this flag to false to run the cleanup process")

	flag.Parse()

	if *cleanupFlag {
		bootStrap()
	} else {
		CleanupCloud()
	}
}
