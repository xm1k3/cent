package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
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

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
