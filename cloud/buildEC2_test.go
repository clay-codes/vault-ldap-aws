package cloud

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	cmdStr := "doormat login"
	cmd := exec.Command("bash", "-c", cmdStr)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}

	cmdStr = "doormat aws export --role $(doormat aws list | tail -n 1 | cut -b 2-)"
	cmd = exec.Command("bash", "-c", cmdStr)
	envVars, _ := cmd.CombinedOutput()
	vars := strings.Split(string(envVars), " && ")

	str := strings.SplitN(strings.TrimPrefix(vars[0], "export "), "=", 2)
	// keep only the value after whitespace in each element
	fmt.Println(str[1])

	//loop that does this for above vars
	for _, declaration := range vars {
		// Removing 'export ' prefix and splitting by '='
		keyValue := strings.SplitN(strings.TrimPrefix(declaration, "export "), "=", 2)
		if len(keyValue) == 2 {
			key, value := keyValue[0], keyValue[1]
			if err := os.Setenv(key, value); err != nil {
				fmt.Printf("Error setting environment variable %s: %v\n", key, err)
			} else {
				fmt.Printf("Set %s=%s\n", key, value)
			}
		}
	}

}

func TestSetRegion(t *testing.T) {
	SetRegion()
}

func TestGetImgID(t *testing.T) {
	// Create a new AWS session with default configuration
	CheckAuth()
	CreateSession()
	err := GetSession().CreateServices()
	if err != nil {
		fmt.Println("Error:", err)
	}
	// Create new EC2 client

	// Call the function under test
	amiID, err := GetImgID()
	os.Setenv("AMI_ID", amiID)
	fmt.Println(os.Getenv("AMI_ID"))
	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, amiID, "AMI ID should not be empty")
}


func TestGetVPC(t *testing.T) {
	CheckAuth()
	CreateSession()
	GetSession().CreateServices("ec2")
	// Call the function under test
	vpcID, err := GetVPC()
	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, vpcID, "VPC ID should not be empty")
}

func TestDescribeVPC(t *testing.T) {
	CheckAuth()
	CreateSession()
	GetSession().CreateServices("ec2")
	vpcs, err := svc.ec2.DescribeVpcs(nil)
	fmt.Println(vpcs)
	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, vpcs, "VPC ID should not be empty")
}

func TestCreateSG(t *testing.T) {
	CheckAuth()
	// Call the function under test
	CreateSession()
	GetSession().CreateServices("ec2")

	sgID, err := CreateSG()

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, sgID, "Security Group ID should not be empty")
}

func TestGetSubnetID(t *testing.T) {
	CheckAuth()
	// Call the function under test
	snID, err := GetSubnetID()
	fmt.Println(os.Getenv("SUBNET_ID"))

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, snID, "Subnet ID should not be empty")
}

// used this to generate a file with encoded user data, then saved encoded string to variable EncodedUserData in buildEC2.go 
func TestEncodeUserData(t *testing.T) {
	// Assume encodeUserData() returns a string
	userData, _ := os.ReadFile("../user-data.yaml")

	encode := base64.StdEncoding.EncodeToString(userData)
	// generate a file with the encoded string
	os.WriteFile("../encoded-user-data.txt", []byte(encode), 0644)

}

func TestBuildEC2(t *testing.T) {
	CheckAuth()

	// Call the function under test
	ec2ID, err := BuildEC2()

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, ec2ID, "EC2 ID should not be empty")
}

func TestGetEC2ID(t *testing.T) {
	CheckAuth()
	// Call the function under test
	CreateSession()
	GetSession().CreateServices("ec2")

	ec2ID, err := GetEC2ID()
	fmt.Println(ec2ID)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, ec2ID, "EC2 ID should not be empty")
}
func TestGetDNS(t *testing.T) {
	CheckAuth()
	// Call the function under test
	CreateSession()
	err := GetSession().CreateServices("ec2")
	if err != nil {
		fmt.Println("Error creating svc", err)
	}
	instanceID, _ := GetEC2ID()
	ipv4dns, err := GetPublicDNS(&instanceID)
	fmt.Println(ipv4dns)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, ipv4dns, "IPv4 DNS should not be empty")
}
