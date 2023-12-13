package cloud

import (
	"os"
	"testing"

	// AWS-specific configurations

	"github.com/stretchr/testify/assert"
)

func TestTerminateEC2Instance(t *testing.T) {
	CheckAuth()
	sess := GetSession().GetAWSSession()

	err := TerminateEC2Instance(sess, os.Getenv("EC2_ID"))
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)
}
func TestDeleteKeyPair(t *testing.T) {
	CheckAuth()
	sess := GetSession().GetAWSSession()

	err := DeleteKeyPair(sess)
	if err != nil {
		t.Fatal(err)
	}
}
func TestDeleSecurityGroup(t *testing.T) {
	CheckAuth()
	sess := GetSession().GetAWSSession()

	err := DeleteSecurityGroup(sess)
	if err != nil {
		t.Fatal(err)
	}

	// call aws describe key pairs to check if key pair is deleted

	assert.NoError(t, err)
}

func TestDetachRoleFromInstanceProfile(t *testing.T) {
	CheckAuth()
	sess := GetSession().GetAWSSession()
	err := DetachRoleFromInstanceProfile(sess)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
}
func TestDeleteInstanceProfile(t *testing.T) {
	CheckAuth()
	sess := GetSession().GetAWSSession()
	err := DeleteInstanceProfile(sess)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
}
func TestDeleteRole(t *testing.T) {
	CheckAuth()
	sess := GetSession().GetAWSSession()
	err := DeleteRole(sess)
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
}
