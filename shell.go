package main

import (
	"fmt"
	"os/exec"
)

func executeShell(cmd string) string {
	cmdout, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return err.Error()
	}
	return fmt.Sprint("\n" + string(cmdout))
}
