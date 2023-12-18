package cloud

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerminateEC2Instance(t *testing.T) {
	CheckAuth()
	err := TerminateEC2Instance()
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)
}
func TestDeleteKeyPair(t *testing.T) {
	CheckAuth()
	err := DeleteKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Remove("../key.pem"); err != nil {
		t.Fatal(err)
	}
}
func TestDeleSecurityGroup(t *testing.T) {
	CheckAuth()
	CreateSession("us-west-2")
	err := GetSession().CreateServices("ec2")
	if err != nil {
		t.Fatal("Error:", err)
	}

	err = DeleteSecurityGroup()
	if err != nil {
		t.Fatal(err)
	}

	// call aws describe key pairs to check if key pair is deleted

	assert.NoError(t, err)
}

func TestDetachRoleFromInstanceProfile(t *testing.T) {
	CheckAuth()
	err := DetachRoleFromInstanceProfile()
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
}
func TestDeleteInstanceProfile(t *testing.T) {
	CheckAuth()
	err := DeleteInstanceProfile()
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
}
func TestDeleteRole(t *testing.T) {
	CheckAuth()
	err := DeleteRole()
	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)
}
