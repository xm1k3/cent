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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"github.com/xm1k3/cent/v2/internal/utils"
	"gopkg.in/yaml.v2"
)

type NucleiTemplate struct {
	Info TemplateInfo `yaml:"info"`
}

type TemplateInfo struct {
	Severity string      `yaml:"severity"`
	Tags     interface{} `yaml:"tags"`
}

type SummaryStats struct {
	Metrics              map[string]int `json:"metrics"`
	SeverityDistribution map[string]int `json:"severity_distribution"`
	Tags                 map[string]int `json:"tags"`
	LastUpdated          string         `json:"last_updated"`
}

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Print summary of nuclei templates by tags and severity",
	Long:  `Print summary of nuclei templates focusing on tags and severity distribution`,
	Run: func(cmd *cobra.Command, args []string) {
		pathFlag, _ := cmd.Flags().GetString("path")
		jsonOutput, _ := cmd.Flags().GetBool("json")
		limit, _ := cmd.Flags().GetInt("limit")
		search, _ := cmd.Flags().GetString("search")

		configDir, err := utils.GetDataDir()
		if err != nil {
			fmt.Printf("Failed to get config directory: %v\n", err)
			os.Exit(1)
		}

		summaryPath := filepath.Join(configDir, "summary.json")

		if _, err := os.Stat(summaryPath); os.IsNotExist(err) {
			fmt.Println("Summary file not found, running update...")
			updateSummary(pathFlag, summaryPath)
		}

		if search != "" {
			searchSummary(summaryPath, search)
			return
		}

		if jsonOutput {
			outputJSON(summaryPath, limit)
		} else {
			outputTable(summaryPath, limit)
		}
	},
}

var summaryUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update summary.json file",
	Long:  `Update the summary.json file with current template statistics`,
	Run: func(cmd *cobra.Command, args []string) {
		pathFlag, _ := cmd.Flags().GetString("path")

		configDir, err := utils.GetDataDir()
		if err != nil {
			fmt.Printf("Failed to get config directory: %v\n", err)
			os.Exit(1)
		}

		summaryPath := filepath.Join(configDir, "summary.json")
		updateSummary(pathFlag, summaryPath)
	},
}

func updateSummary(path string, summaryPath string) {
	stats := analyzeTemplates(path)
	stats.LastUpdated = time.Now().Format("2006-01-02 15:04:05")

	jsonData, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	err = os.WriteFile(summaryPath, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing summary file: %v\n", err)
		return
	}

	fmt.Printf("Summary updated and saved to: %s\n", summaryPath)
}

func analyzeTemplates(path string) SummaryStats {
	stats := SummaryStats{
		Metrics:              make(map[string]int),
		SeverityDistribution: make(map[string]int),
		Tags:                 make(map[string]int),
	}

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml") {
			stats.Metrics["total_templates"]++

			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				stats.Metrics["invalid_templates"]++
				return err
			}

			var template NucleiTemplate
			err = yaml.Unmarshal(fileContent, &template)
			if err != nil {
				stats.Metrics["invalid_templates"]++
				os.Remove(filePath)
				return err
			}

			if template.Info.Severity != "" {
				severity := strings.ToUpper(template.Info.Severity)
				stats.SeverityDistribution[severity]++
			}

			tags, ok := template.Info.Tags.([]interface{})
			if ok {
				for _, tag := range tags {
					if tagStr, ok := tag.(string); ok {
						tagParts := strings.Split(tagStr, ",")
						for _, tagPart := range tagParts {
							tagPart = strings.TrimSpace(tagPart)
							if tagPart != "" {
								stats.Tags[tagPart]++

								if strings.Contains(strings.ToLower(tagPart), "cve") {
									stats.Metrics["cve_templates"]++
								}
							}
						}
					}
				}
			} else if tagStr, ok := template.Info.Tags.(string); ok {
				tagParts := strings.Split(tagStr, ",")
				for _, tagPart := range tagParts {
					tagPart = strings.TrimSpace(tagPart)
					if tagPart != "" {
						stats.Tags[tagPart]++

						if strings.Contains(strings.ToLower(tagPart), "cve") {
							stats.Metrics["cve_templates"]++
						}
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}

	stats.Metrics["valid_templates"] = stats.Metrics["total_templates"] - stats.Metrics["invalid_templates"]

	return stats
}

func loadSummary(summaryPath string) SummaryStats {
	fileContent, err := os.ReadFile(summaryPath)
	if err != nil {
		fmt.Printf("Error reading summary file: %v\n", err)
		os.Exit(1)
	}

	var stats SummaryStats
	err = json.Unmarshal(fileContent, &stats)
	if err != nil {
		fmt.Printf("Error parsing summary file: %v\n", err)
		os.Exit(1)
	}

	return stats
}

func searchSummary(summaryPath string, search string) {
	stats := loadSummary(summaryPath)
	search = strings.ToLower(search)

	fmt.Printf("Search results for: %s\n\n", search)

	found := false

	for key, value := range stats.Metrics {
		if strings.Contains(strings.ToLower(key), search) {
			fmt.Printf("Metric - %s: %d\n", key, value)
			found = true
		}
	}

	for key, value := range stats.SeverityDistribution {
		if strings.Contains(strings.ToLower(key), search) {
			fmt.Printf("Severity - %s: %d\n", key, value)
			found = true
		}
	}

	for key, value := range stats.Tags {
		if strings.Contains(strings.ToLower(key), search) {
			fmt.Printf("Tag - %s: %d\n", key, value)
			found = true
		}
	}

	if !found {
		fmt.Println("No results found")
	}
}

func outputJSON(summaryPath string, limit int) {
	stats := loadSummary(summaryPath)

	if limit > 0 {
		limitedTags := make(map[string]int)
		type tagCount struct {
			tag   string
			count int
		}
		var tagCounts []tagCount
		for tag, count := range stats.Tags {
			tagCounts = append(tagCounts, tagCount{tag, count})
		}
		sort.Slice(tagCounts, func(i, j int) bool {
			return tagCounts[i].count > tagCounts[j].count
		})

		for i := 0; i < len(tagCounts) && i < limit; i++ {
			limitedTags[tagCounts[i].tag] = tagCounts[i].count
		}
		stats.Tags = limitedTags
	}

	jsonData, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

func outputTable(summaryPath string, limit int) {
	stats := loadSummary(summaryPath)

	fmt.Println(color.CyanString("\n=== NUCLEI TEMPLATES SUMMARY ===\n"))

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Metric", "Count"})
	t.AppendRows([]table.Row{
		{"Total Templates", stats.Metrics["total_templates"]},
		{"CVE Templates", stats.Metrics["cve_templates"]},
		{"Invalid Templates", stats.Metrics["invalid_templates"]},
		{"Valid Templates", stats.Metrics["valid_templates"]},
	})
	t.Render()

	if len(stats.SeverityDistribution) > 0 {
		fmt.Println(color.YellowString("\n=== SEVERITY DISTRIBUTION ==="))
		t2 := table.NewWriter()
		t2.SetOutputMirror(os.Stdout)
		t2.AppendHeader(table.Row{"Severity", "Count"})

		severityOrder := []string{"CRITICAL", "HIGH", "MEDIUM", "LOW", "INFO"}
		for _, severity := range severityOrder {
			if count, exists := stats.SeverityDistribution[severity]; exists {
				t2.AppendRow(table.Row{severity, count})
			}
		}
		t2.Render()
	}

	if len(stats.Tags) > 0 {
		fmt.Println(color.YellowString("\n=== TOP TAGS ==="))
		t3 := table.NewWriter()
		t3.SetOutputMirror(os.Stdout)
		t3.AppendHeader(table.Row{"Tag", "Count"})

		type tagCount struct {
			tag   string
			count int
		}
		var tagCounts []tagCount
		for tag, count := range stats.Tags {
			tagCounts = append(tagCounts, tagCount{tag, count})
		}
		sort.Slice(tagCounts, func(i, j int) bool {
			return tagCounts[i].count > tagCounts[j].count
		})

		displayLimit := len(tagCounts)
		if limit > 0 && limit < displayLimit {
			displayLimit = limit
		}

		for i := 0; i < displayLimit; i++ {
			t3.AppendRow(table.Row{tagCounts[i].tag, tagCounts[i].count})
		}
		t3.Render()
	}

	fmt.Printf(color.CyanString("\nLast updated: %s\n"), stats.LastUpdated)
}

func init() {
	rootCmd.AddCommand(summaryCmd)
	summaryCmd.AddCommand(summaryUpdateCmd)

	summaryCmd.Flags().StringP("path", "p", "cent-nuclei-templates", "Root path to save the templates")
	summaryCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	summaryCmd.Flags().IntP("limit", "l", 25, "Limit number of tags to display")
	summaryCmd.Flags().StringP("search", "s", "", "Search for specific data in summary")

	summaryUpdateCmd.Flags().StringP("path", "p", "cent-nuclei-templates", "Root path to save the templates")
}
