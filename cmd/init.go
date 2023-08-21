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

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/xm1k3/cent/internal/utils"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Cent init configuration file",
	Long:  "This command will automatically download .cent.yaml from repo and copy it to $HOME/.cent.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		var fileUrl string
		linkFlag, _ := cmd.Flags().GetString("url")
		home, _ := homedir.Dir()
		overwrite, _ := cmd.Flags().GetBool("overwrite")

		if _, err := os.Stat(home + "/.cent.yaml"); !os.IsNotExist(err) {
			if !overwrite {
				log.Fatal("Cent config file already exists, if you want to overwrite it use the --overwrite flag ")
			}
		}

		if linkFlag == "" {
			fileUrl = "https://raw.githubusercontent.com/xm1k3/cent/main/.cent.yaml"
		} else {
			fileUrl = linkFlag
		}
		err := utils.DownloadFile(home+"/.cent.yaml", fileUrl)
		if err != nil {
			panic(err)
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("cent configured correctly, you can find the configuration file here: " + home + "/.cent.yaml")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("url", "u", "", "Url from which you can download the configurations for .cent.yaml")
	initCmd.Flags().BoolP("overwrite", "o", false, "If the cent file exists overwrite it")
}
