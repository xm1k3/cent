/*
Copyright © 2023 xm1k3

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xm1k3/cent/internal/utils"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Cent init configuration file",
	Long:  "This command will automatically download .cent.yaml from repo and copy it to .config/cent/.cent.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		var fileUrl string
		linkFlag, _ := cmd.Flags().GetString("url")
		overwrite, _ := cmd.Flags().GetBool("overwrite")

		configDir, err := utils.GetDataDir()
		if err != nil {
			log.Fatalf("Failed to get config directory: %v", err)
		}

		if err := os.MkdirAll(configDir, 0755); err != nil {
			log.Fatalf("Failed to create config directory %s: %v", configDir, err)
		}

		configFilePath := filepath.Join(configDir, ".cent.yaml")

		if _, err := os.Stat(configFilePath); !os.IsNotExist(err) {
			if !overwrite {
				log.Fatal("Cent config file already exists, if you want to overwrite it use the --overwrite flag")
			}
		}

		if linkFlag == "" {
			fileUrl = "https://raw.githubusercontent.com/xm1k3/cent/main/.cent.yaml"
		} else {
			fileUrl = linkFlag
		}

		err = utils.DownloadFile(configFilePath, fileUrl)
		if err != nil {
			log.Fatalf("Failed to download config file: %v", err)
		}

		fmt.Printf("Cent configured correctly, you can find the configuration file here: %s\n", configFilePath)
	},
}

var initCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if .cent.yaml configuration file exists",
	Long:  "Check if the .cent.yaml configuration file exists in the config directory",
	Run: func(cmd *cobra.Command, args []string) {
		configDir, err := utils.GetDataDir()
		if err != nil {
			log.Fatalf("Failed to get config directory: %v", err)
		}

		configFilePath := filepath.Join(configDir, ".cent.yaml")

		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			fmt.Printf(color.RedString("Configuration file not found at: %s\n"), configFilePath)
			fmt.Printf(color.YellowString("Run 'cent init' to download the configuration file\n"))
			os.Exit(1)
		} else if err != nil {
			fmt.Printf(color.RedString("Error checking configuration file: %v\n"), err)
			os.Exit(1)
		} else {
			fmt.Printf(color.GreenString("Configuration file found at: %s\n"), configFilePath)

			if fileInfo, err := os.Stat(configFilePath); err == nil {
				fmt.Printf(color.CyanString("File size: %d bytes\n"), fileInfo.Size())
				fmt.Printf(color.CyanString("Last modified: %s\n"), fileInfo.ModTime().Format("2006-01-02 15:04:05"))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.AddCommand(initCheckCmd)

	initCmd.Flags().StringP("url", "u", "", "Url from which you can download the configurations for .cent.yaml")
	initCmd.Flags().BoolP("overwrite", "o", false, "If the cent file exists overwrite it")
}
