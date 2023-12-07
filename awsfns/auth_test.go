package awsfns

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
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
