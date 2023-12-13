package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/clay-codes/aws-ldap/cloud"
)

func bootStrap(sess *session.Session) {
	imgID, err := cloud.GetImgID(sess)
	if err != nil {
		log.Fatal(err)
	}

	str, err := cloud.CreateKP(sess)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	vpcID, err := cloud.GetVPC(sess)
	if err != nil {
		log.Fatal(err)
	}

	sgid, err := cloud.CreateSG(sess, vpcID, []int64{22, 8200, 8201})
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("SGID", sgid)

	err = cloud.CreateInstProf(sess)
	if err != nil {
		log.Fatal(err)
	}

	snID, err := cloud.GetSubnetID(sess, *vpcID)
	if err != nil {
		log.Fatal(err)
	}

	ec2ID, err := cloud.BuildEC2(sess, []string{sgid}, imgID, snID)
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("EC2_ID", ec2ID)
	fmt.Printf("EC2ID:%s\n", ec2ID)
}

func cleanupCloud(sess *session.Session, instanceID string, sgID string) error {
	if err := os.Remove("key.pem"); err != nil {
		return fmt.Errorf("error deleting key.pem file: %v", err)
	}
	if err := cloud.TerminateEC2Instance(sess, os.Getenv("EC2_ID")); err != nil {
		return fmt.Errorf("error terminating EC2 instance: %v", err)
	}
	if err := cloud.DeleteKeyPair(sess); err != nil {
		return fmt.Errorf("error deleting key pair: %v", err)
	}

	if err := cloud.DeleteSecurityGroup(sess); err != nil {
		return fmt.Errorf("error deleting security group: %v", err)
	}
	if err := cloud.DetachRoleFromInstanceProfile(sess); err != nil {
		return fmt.Errorf("error detaching role from instance profile: %v", err)
	}
	if err := cloud.DeleteInstanceProfile(sess); err != nil {
		return fmt.Errorf("error deleting instance profile: %v", err)
	}
	if err := cloud.DeleteRole(sess); err != nil {
		return fmt.Errorf("error deleting role: %v", err)
	}
	return nil
}

// run the main function
func main() {
	// Initialize a session in us-west-2 that the SDK will use to load credentials
	// from the shared credentials file ~/.aws/credentials.

	//set flag to run cleanup
	cloud.CheckAuth()

	cleanupFlag := flag.Bool("cleanup", false, "Set this flag to run the cleanup process")

	flag.Parse()

	err := cloud.CreateSession("us-west-2")
	if err != nil {
		log.Fatal(err)
	}

	sess := cloud.GetSession().GetAWSSession()

	if *cleanupFlag {
		bootStrap(sess)
	} else {
		cleanupCloud(sess, os.Getenv("EC2_ID"), os.Getenv("SGID"))
	}
}
