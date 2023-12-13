/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
	"github.com/xm1k3/cent/pkg/jobs"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate templates",
	Long: `The validate command is a part of the application's functionality to validate templates. 
When executed, it scans a specified folder for YAML files. Each YAML file is checked for validity. 
If a template is found to be invalid, it is deleted from the folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")

		var wg sync.WaitGroup
		filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
			}
			if info.IsDir() {
			} else if filepath.Ext(filePath) == ".yaml" {
				wg.Add(1)
				go func(filePath string) {
					defer wg.Done()

					data, err := readFile(filePath)
					if err != nil {
						return
					}

					isValid, err := jobs.ValidateTemplate(data)
					if err != nil {
						fmt.Println(filePath, err)
					}

					if !isValid {
						os.Remove(filePath)
					}
				}(filePath)
			}
			return nil
		})

		wg.Wait()
	},
}

func readFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("path", "p", "cent-nuclei-templates", "Root path to save the templates")
}
