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
	fmt.Print("Would you like to run cleanup? ")
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
	key, err := cloud.CreateKP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nkey pair created            %s", key)

	sgid, err := cloud.CreateSG()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nsecurity group created      %s", sgid)
	err = cloud.CreateInstProf()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nrole created                ec2-admin-role-custom")
	fmt.Println("instance profile created    ec2-InstProf-custom")
	// wait for instance profile to be created sometimes necessary to avoid not found error
	time.Sleep(10 * time.Second)

	pubDNS, err := cloud.BuildEC2()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nusername: Administrator")
	fmt.Println("password: admin")
	fmt.Println("\n\nEnvironment nearly ready. Use this command to connect via openssh in a few moments: ")
	fmt.Printf("\nssh -i key.pem -o StrictHostKeyChecking=no Administrator@%s\n", pubDNS)
	fmt.Println("\n\nAD server installed with forest dn of vaultest.com \nUse the following command to test the connection via ldapsearch: ")
	fmt.Printf("\nldapsearch -x -H ldap://%s:389 -D \"cn=admin,dc=vaultest,dc=com\" -w admin -b \"dc=vaultest,dc=com\" -s sub \"(objectclass=*)\"\n\n", pubDNS)
	fmt.Println("\nCan also verify forest (root dn) exists by running the following on the server itself: ")
	fmt.Println("\n> powershell")
	fmt.Println("> Import-Module C:\\Windows\\system32\\WindowsPowerShell\\v1.0\\Modules\\ActiveDirectory\\ActiveDirectory.psd1")
	fmt.Println("> Get-ADForest")
}

func CleanupCloud() {
	if err := os.Remove("key.pem"); err != nil {
		fmt.Printf("key.pem file may not exist: %v\n", err)
	}
	if err := cloud.TerminateEC2Instance(); err != nil {
		fmt.Printf("instance may not have been created: %v\n", err)
	}
	if err := cloud.DeleteKeyPair(); err != nil {
		fmt.Printf("key pair may not exist: %v\n", err)
	}

	if err := cloud.DetachPolicyFromRole(); err != nil {
		fmt.Printf("policy may not have been created: %v\n", err)
	}

	if err := cloud.DetachRoleFromInstanceProfile(); err != nil {
		fmt.Printf("error detaching role from instance profile: %v\n", err)
	}
	if err := cloud.DeleteInstanceProfile(); err != nil {
		fmt.Printf("error deleting instance profile: %v\n", err)
	}
	if err := cloud.DeleteRole(); err != nil {
		fmt.Printf("error deleting role: %v\n", err)
	}
	if err := cloud.DeleteSecurityGroup(); err != nil {
		fmt.Printf("error deleting security group: %v\n", err)
	}
}

func main() {
	if runCleanup {
		CleanupCloud()
	} else {
		bootStrap()
	}
}
