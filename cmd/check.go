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
	"log"

	"github.com/spf13/cobra"
	"github.com/xm1k3/cent/pkg/jobs"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if templates repo are still available",
	Long:  `Check if templates repo are still available`,
	Run: func(cmd *cobra.Command, args []string) {
		removeFlag, _ := cmd.Flags().GetBool("remove")
		err := jobs.CheckConfig(cfgFile, removeFlag)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().BoolP("remove", "r", false, "Remove from .cent.yaml urls that are no longer accessible or are currently private")

}
