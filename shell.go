package main

import (
	"fmt"
	"os/exec"
)

func executeShell(program string) (string, error) {
	cmd := exec.Command("sh", "-c", program)
	cmdout, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("executed program failed on exit %s :%v", program, err)
	}
	return fmt.Sprint("\n" + string(cmdout)), nil
}

func checkProgramExists(program string) bool {
	_, err := exec.LookPath(program)
	if err != nil {
		return true
	}
	return false
}
