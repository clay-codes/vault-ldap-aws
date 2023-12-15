package main

import (
	"testing"

	// AWS-specific configurations

	"github.com/clay-codes/aws-ldap/cloud"
)

func TestCleanupCloud(t *testing.T) {
	cloud.CheckAuth()
	//ec2ID := "i-0716c3d91881333a1"

	CleanupCloud()

}

func TestBootStrap(t *testing.T) {
	cloud.CheckAuth()
	bootStrap()
}
