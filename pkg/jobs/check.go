package jobs

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ExcludeDirs        []string    `yaml:"exclude-dirs"`
	ExcludeFiles       []string    `yaml:"exclude-files"`
	CommunityTemplates []yaml.Node `yaml:"community-templates"`
}

func configEntryURL(node yaml.Node) string {
	if node.Kind == yaml.ScalarNode {
		return node.Value
	}
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content)-1; i += 2 {
			if node.Content[i].Value == "url" {
				return node.Content[i+1].Value
			}
		}
	}
	return ""
}

func CheckConfig(configPath string, removeFlag bool) error {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, node := range config.CommunityTemplates {
		url := configEntryURL(node)
		if url == "" {
			continue
		}

		if removeFlag && strings.Contains(url, "gist.github.com") {
			fmt.Println(color.YellowString("[INFO] Ignoring and removing gist.github.com URL: %s", url))
			RemoveURLFromConfig(configPath, url)
			continue
		}

		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Error on GET request for %s: %s\n", url, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				fmt.Println(color.GreenString("[SUCCESS] URL %s Status code: %d", url, resp.StatusCode))
				return
			} else if resp.StatusCode == http.StatusNotFound {
				fmt.Println(color.RedString("[ERR] URL %s  Status code: %d", url, resp.StatusCode))
				if removeFlag {
					RemoveURLFromConfig(configPath, url)
				}
			} else {
				fmt.Println(color.RedString("[ERR] URL %s  Status code: %d", url, resp.StatusCode))
			}
		}(url)
	}

	wg.Wait()

	return nil
}

func RemoveURLFromConfig(configPath, url string) error {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return err
	}

	var filtered []yaml.Node
	for _, node := range config.CommunityTemplates {
		if configEntryURL(node) != url {
			filtered = append(filtered, node)
		}
	}

	updatedConfig := Config{
		ExcludeDirs:        config.ExcludeDirs,
		ExcludeFiles:       config.ExcludeFiles,
		CommunityTemplates: filtered,
	}

	updatedYAML, err := yaml.Marshal(&updatedConfig)
	if err != nil {
		return err
	}

	updatedYAMLString := string(updatedYAML)
	updatedYAMLString = strings.ReplaceAll(updatedYAMLString, "  ", " ")

	err = os.WriteFile(configPath, []byte(updatedYAMLString), 0644)
	if err != nil {
		return err
	}

	return nil
}
