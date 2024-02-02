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
	CreateSession()
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

func TestBootPrint(t *testing.T) {
	key:= "key"
	sgid:= "sgid"
	pubDNS:= "pubDNS"
	fmt.Printf("\nkey pair created            %s", key)
	fmt.Printf("\nsecurity group created      %s", sgid)
	fmt.Println("\nrole created                ec2-admin-role-custom")
	fmt.Println("instance profile created    ec2-InstProf-custom")
	fmt.Println("\nusername            Administrator")
	fmt.Println("password            admin")
	fmt.Println("forest (root) dn    DC=vaultest,DC=com")
	fmt.Println("\n\nEnvironment nearly ready. Server will need an additional few minutes to bootstrap AD even after connection established. ")
	fmt.Println("\nRun this to connect: ")
	fmt.Printf("\nssh -i key.pem -o StrictHostKeyChecking=no Administrator@%s\n", pubDNS)
	fmt.Println("\n\nUse the following command to test the connection: ")
	fmt.Printf("\nldapsearch -x -H ldap://%s:389 -D \"cn=admin,dc=vaultest,dc=com\" -w admin -b \"dc=vaultest,dc=com\" -s sub \"(objectclass=*)\"\n\n", pubDNS)
	fmt.Println("\nCan also verify forest (root dn) exists on the server via: ") 
	fmt.Println("\n> powershell")
	fmt.Println("> Import-Module C:\\Windows\\system32\\WindowsPowerShell\\v1.0\\Modules\\ActiveDirectory\\ActiveDirectory.psd1")
	fmt.Println("> Get-ADForest")
}

func TestCleanupPrint(t *testing.T) {
	instanceID:= "i-1234567890abcdef0"
	sgID:= "sg-1234567890abcdef0"
	fmt.Println("\nEC2 instance terminated                       ", instanceID)
	fmt.Println("Key pair deleted                               vault-EC2-kp")
	fmt.Println("AWS role detatched from custom role            AmazonSSMAutomationRole")
	fmt.Println("Custom role detached from instance profile     ec2-admin-role-custom")
	fmt.Println("Instance profile deleted                       ec2-InstProf-custom")
	fmt.Println("Custom role deleted                            ec2-admin-role-custom")
	fmt.Println("Security group deleted                        ", sgID)

	
	
}