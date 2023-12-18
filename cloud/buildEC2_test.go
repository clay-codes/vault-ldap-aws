package cloud

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	// AWS-specific configurations

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
	//fmt.Println(len(vars))
	// fmt.Println(vars[0])
	// fmt.Print(strings.TrimPrefix(vars[0], "export "))
	// fmt.Print(strings.TrimPrefix(vars[0], "export "))
	//split string at first '=' and trim prefix string, including whitespace

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

	//vars

	//log.Printf("Output: %s", output)
	// get first two strings from 'output'
	// words := strings.Fields(string(output))
	// log.Print(words[:2])
	// fmt.Println(words[:2])
	// pattern := `AWS_SESSION_TOKEN=[^\s]+`
	// re:= regexp.MustCompile(pattern)
	// match := re.FindString(string(output))

	// re = regexp.MustCompile(`=(.*)`)
	// match = re.FindStringSubmatch(match)[1]

	// if match != "" {
	// 	fmt.Println(match)
	// } else {
	// 	fmt.Println("No match found")
	// }

	//os.Setenv("AWS_SESSION_TOKEN", match)

	// for _, declaration := range vars {
	// 	// Removing 'export ' prefix and splitting by '='
	// 	keyValue := strings.Split(strings.TrimPrefix(declaration, "export "), "=")
	// 	if len(keyValue) == 2 {
	// 		key, value := keyValue[0], keyValue[1]
	// 		if err := os.Setenv(key, value); err != nil {
	// 			fmt.Printf("Error setting environment variable %s: %v\n", key, err)
	// 		} else {
	// 			fmt.Printf("Set %s=%s\n", key, value)
	// 		}
	// 	}
	// }

	//}

	// Additional assertions can be added here as needed
}

func TestGetImgID(t *testing.T) {
	// Create a new AWS session with default configuration
	CheckAuth()
	CreateSession("us-west-2")
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
	// Call the function under test
	vpcID, err := GetVPC()
	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, vpcID, "VPC ID should not be empty")
}

func TestCreateSG(t *testing.T) {
	CheckAuth()
	// Call the function under test
	CreateSession("us-west-2")
	err := GetSession().CreateServices("ec2")
	if err != nil {
		fmt.Println("Error:", err)
	}
	sgID, err := CreateSG([]int64{22, 8200, 8201})

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

func TestBuildEC2(t *testing.T) {
	CheckAuth()

	// Call the function under test
	ec2ID, err := BuildEC2()

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, ec2ID, "EC2 ID should not be empty")
}
