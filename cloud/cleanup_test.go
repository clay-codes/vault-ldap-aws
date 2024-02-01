package cloud

import (
	"fmt"
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

func TestPrint(t *testing.T) {
	key:= "key"
	sgid:= "sgid"
	pubDNS:= "pubDNS"
	fmt.Printf("\nkey pair created            %s", key)
	fmt.Printf("\nsecurity group created      %s", sgid)
	fmt.Println("\nrole created                ec2-admin-role-custom")
	fmt.Println("instance profile created    ec2-InstProf-custom")
	fmt.Println("\nusername: Administrator")
	fmt.Println("password: admin")
	fmt.Println("\n\nEnvironment nearly ready. Use this command to connect via openssh in a few moments: ")
	fmt.Printf("\nssh -i key.pem -o StrictHostKeyChecking=no Administrator@%s\n", pubDNS)
	fmt.Println("\n\nAD server installed with forest dn of vaultest.com. Use the following command to test the connection: ")
	fmt.Printf("\nldapsearch -x -H ldap://%s:389 -D \"cn=admin,dc=vaultest,dc=com\" -w admin -b \"dc=vaultest,dc=com\" -s sub \"(objectclass=*)\"\n\n", pubDNS)
	fmt.Println("\nCan also verify forest (root dn) exists on the server via: ") 
	fmt.Println("\n> powershell")
	fmt.Println("> Import-Module C:\\Windows\\system32\\WindowsPowerShell\\v1.0\\Modules\\ActiveDirectory\\ActiveDirectory.psd1")
	fmt.Println("> Get-ADForest")
}