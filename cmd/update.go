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
	"github.com/xm1k3/cent/pkg/jobs"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update your repository",
	Long:  `This command helps you update your folder with templates by deleting unnecessary folders and files without having to do multiples git clones.`,
	Run: func(cmd *cobra.Command, args []string) {
		pathFlag, _ := cmd.Flags().GetString("path")
		directories, _ := cmd.Flags().GetBool("directories")
		files, _ := cmd.Flags().GetBool("files")

		configPath := cfgFile
		if configPath == "" {
			configDir, err := utils.GetDataDir()
			if err != nil {
				log.Fatalf("Failed to get config directory: %v", err)
			}
			configPath = filepath.Join(configDir, ".cent.yaml")
		}

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Println(`Run ` + color.YellowString("cent init") + ` to automatically download ` +
				color.HiCyanString(".cent.yaml") + ` from repo and copy it to ` +
				color.HiCyanString(".config/cent/.cent.yaml"))
			return
		}

		if pathFlag != "" {
			if directories || files {
				jobs.UpdateRepo(pathFlag, directories, files, true)
			} else {
				fmt.Println(color.YellowString("[!] directory or file flag required"))
			}
		} else {
			fmt.Println(color.YellowString("[!] Set folder path flag (example: $HOME/cent)"))
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolP("directories", "d", false, "If true remove unnecessary folders from updated .cent.yaml")
	updateCmd.Flags().BoolP("files", "f", false, "If true remove unnecessary files from updated .cent.yaml")
	updateCmd.Flags().StringP("path", "p", "", "Path to folder with nuclei templates")
}
