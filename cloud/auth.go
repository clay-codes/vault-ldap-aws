package cloud

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Auth() (string, error) {
	cmd := exec.Command("bash", "-c", "doormat login")
	if err := cmd.Run(); err != nil {
		return "doormat CLI not installed: ", err
	}

	// getting export statements from doormat aws export
	// gives all commands in one string
	cmd = exec.Command("bash", "-c", "doormat aws export --role $(doormat aws list | tail -n 1 | cut -b 2-)")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "doormat CLI installed, but now issue with AWS creds: ", err
	}

	// Splitting output by ' && ' to get individual variable declaration commands into an array
	// example: [export AWS_ACCESS_KEY_ID=ASIA..., export AWS_SECRET_ACCESS_KEY=..., ...]
	varDeclarations := strings.Split(string(output), " && ")

	// this loop sets the environment variables for the current process using os.Setenv
	for _, declaration := range varDeclarations {
		// for each element, creates slice with 'export ' prefix removed and split by first occurence of '='
		// result: [AWS_ACCESS_KEY_ID, ASIAXYZ...]
		// then sets environment variable using os.Setenv with key and value
		keyValue := strings.SplitN(strings.TrimPrefix(declaration, "export "), "=", 2)
		if len(keyValue) == 2 {
			key, value := keyValue[0], keyValue[1]
			if err := os.Setenv(key, value); err != nil {
				fmt.Printf("Error setting environment variable %s: %v\n", key, err)
			}
		}
	}
	return "auth succeeded", nil
}
