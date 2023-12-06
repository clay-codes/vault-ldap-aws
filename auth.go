package main

import (
	"os/exec"
)

func Auth() (string, error) {
	cmdStr := "doormat login && eval `doormat aws export --role $(doormat aws list | tail -n 1 | cut -b 2-)`"
	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
