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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Info struct {
	Tags string `yaml:"tags"`
}

type YamlFile struct {
	Info Info `yaml:"info"`
}

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Print summary",
	Long:  `Print summary`,
	Run: func(cmd *cobra.Command, args []string) {
		pathFlag, _ := cmd.Flags().GetString("path")

		yamlFilesWithCVE := 0
		yamlFilesWithoutCVE := 0
		totalYamlFiles := 0

		err := filepath.Walk(pathFlag, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml") {
				totalYamlFiles++

				fileContent, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				var yamlFile YamlFile
				err = yaml.Unmarshal(fileContent, &yamlFile)
				if err != nil {
					os.Remove(path)
				}

				tags := strings.Split(yamlFile.Info.Tags, ",")
				tagFound := false

				for _, tag := range tags {
					if strings.Contains(strings.ToLower(tag), "cve") {
						tagFound = true
						break
					}
				}

				if tagFound {
					yamlFilesWithCVE++
				} else {
					yamlFilesWithoutCVE++
				}
			}

			return nil
		})

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Templates type", "Templates count"})
		t.AppendRows([]table.Row{
			{"CVE Templates", yamlFilesWithCVE},
			{"Other Vulnerability Templates", yamlFilesWithoutCVE},
			{"Total Templates", totalYamlFiles},
		})
		t.Render()
	},
}

func init() {
	rootCmd.AddCommand(summaryCmd)

	summaryCmd.Flags().StringP("path", "p", "cent-nuclei-templates", "Root path to save the templates")
}
