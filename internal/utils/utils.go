package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func RunCommand(command string, background bool) {
	cmd := exec.Command("bash", "-c", command)
	if background {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		var execOut bytes.Buffer
		var execErr bytes.Buffer
		cmd.Stdout = &execOut
		cmd.Stderr = &execErr
	}
	err := cmd.Run()
	if err != nil {
		if background {
			fmt.Println("Error running shell command: ", command, "  => ", err.Error())
		}
	}
}
