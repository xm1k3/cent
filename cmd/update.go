/*
Copyright © 2021 xm1k3

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

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xm1k3/cent/pkg/jobs"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update your repository",
	Long:  `This command helps you update your folder with templates by deleting unnecessary folders and files without having to do multiples git clone.`,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := rootCmd.PersistentFlags().GetString("path")
		directories, _ := cmd.Flags().GetBool("directories")
		files, _ := cmd.Flags().GetBool("files")
		if path != "" {
			if directories || files {
				jobs.UpdateRepo(path, directories, files)
			} else {
				fmt.Println(color.YellowString("[!] See avaiable flags"))
			}
		} else {
			fmt.Println(color.YellowString("[!] Set path flag"))
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolP("directories", "d", false, "Remove unnecessary folders from updated $HOME/.cent.yaml")
	updateCmd.Flags().BoolP("files", "f", false, "Remove unnecessary files from updated $HOME/.cent.yaml")

}
