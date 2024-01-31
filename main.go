package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/clay-codes/aws-ldap/cloud"
)

var runCleanup bool

func init() {

	// prompt user if they want to run cleanup
	fmt.Print("Would you like to run cleanup? (yes/no): ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	runCleanup = strings.ToLower(response) == "yes" || strings.ToLower(response) == "y"

	// authenticate with AWS
	cloud.CheckAuth()

	// creating a session
	if err := cloud.CreateSession("us-west-2"); err != nil {
		log.Fatal(err)
	}

	// creating needed services from session
	if err := cloud.GetSession().CreateServices(); err != nil {
		log.Fatal(err)
	}
}

// build environment
func bootStrap() {
	str, err := cloud.CreateKP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)

	_, err = cloud.CreateSG()
	if err != nil {
		log.Fatal(err)
	}

	err = cloud.CreateInstProf()
	if err != nil {
		log.Fatal(err)
	}
	// wait for instance profile to be created sometimes necessary to avoid not found error
	time.Sleep(10 * time.Second)

	pubDNS, err := cloud.BuildEC2()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Environment nearly ready. Use this command to connect via openssh in a few moments: ")
	fmt.Printf("ssh -i key.pem -o StrictHostKeyChecking=no Administrator@%s", pubDNS)
	fmt.Println()
	fmt.Println()
	fmt.Println("Root Username: Administrator")
	fmt.Println("Password: admin")
	fmt.Println("Test DN of vaultest.com has been added to the LDAP server. Use the following command to test the connection: ")
	fmt.Printf("ldapsearch -x -H ldap://%s:389 -D \"cn=admin,dc=vaultest,dc=com\" -w admin -b \"dc=vaultest,dc=com\" -s sub \"(objectclass=*)\"\n", pubDNS)
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

	if err := cloud.DetachPolicyFromRole(); err != nil {
		fmt.Printf("error detaching policy from role: %v", err)
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

func main() {
	if runCleanup {
		CleanupCloud()
	} else {
		bootStrap()
	}
}
