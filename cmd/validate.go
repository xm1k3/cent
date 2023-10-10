/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xm1k3/cent/pkg/jobs"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate templates",
	Long:  `Validate templates`,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := rootCmd.Flags().GetString("path")

		filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
			} else if filepath.Ext(filePath) == ".yaml" {
				data, err := readFile(filePath)
				if err != nil {
					return err
				}

				isValid, err := jobs.ValidateTemplate(data)
				if err != nil {
					fmt.Println(filePath, err)
				}

				if !isValid {
					os.Remove(filePath)
				}
			}
			return nil
		})

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
}
