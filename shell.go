package main

import (
	"fmt"
	"os/exec"
)

func executeShell(cmd string) string {
	if cmdout, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return err.Error()
	} else {
		return fmt.Sprint("\n" + string(cmdout))
	}
}
