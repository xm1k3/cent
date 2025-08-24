/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
	"github.com/xm1k3/cent/v2/pkg/jobs"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate templates",
	Long: `The validate command is a part of the application's functionality to validate templates. 
When executed, it scans a specified folder for YAML files. Each YAML file is checked for validity. 
If a template is found to be invalid, it is deleted from the folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		maxWorkers := 1000

		path, _ := cmd.Flags().GetString("path")
		var wg sync.WaitGroup

		
		workerPool := make(chan struct{}, maxWorkers)

		err := processFiles(path, &wg, workerPool)
		if err != nil {
			log.Fatal(err)
		}

		wg.Wait()
	},
}

// Walk through files and process them
func processFiles(path string, wg *sync.WaitGroup, workerPool chan struct{}) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(filePath) == ".yaml" {
			wg.Add(1)
			processFile(filePath, wg, workerPool)
		}

		return nil
	})
}

// Process an individual file
func processFile(filePath string, wg *sync.WaitGroup, workerPool chan struct{}) {
	// Acquire a worker slot from the pool
	workerPool <- struct{}{}

	go func(filePath string) {
		defer func() {
			wg.Done()
			<-workerPool // Release the worker slot
		}()

		if err := validateTemplate(filePath); err != nil {
			fmt.Println(filePath, err)
		}
	}(filePath)
}

// Validate a single template file
func validateTemplate(filePath string) error {
	data, err := readFile(filePath)
	if err != nil {
		return err
	}

	isValid, err := jobs.ValidateTemplate(data)
	if err != nil {
		return err
	}

	if !isValid {
		return os.Remove(filePath)
	}

	return nil
}

func readFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("path", "p", "cent-nuclei-templates", "Root path to save the templates")

}
