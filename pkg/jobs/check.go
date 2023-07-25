package jobs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/xm1k3/cent/internal/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ExcludeDirs        []string `yaml:"exclude-dirs"`
	ExcludeFiles       []string `yaml:"exclude-files"`
	CommunityTemplates []string `yaml:"community-templates"`
}

func CheckConfig(configPath string, removeFlag bool) error {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, url := range config.CommunityTemplates {
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
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return err
	}

	updatedConfig := Config{
		ExcludeDirs:        config.ExcludeDirs,
		ExcludeFiles:       config.ExcludeFiles,
		CommunityTemplates: utils.RemoveStringFromSlice(config.CommunityTemplates, url),
	}

	updatedYAML, err := yaml.Marshal(&updatedConfig)
	if err != nil {
		return err
	}

	updatedYAMLString := string(updatedYAML)
	updatedYAMLString = strings.ReplaceAll(updatedYAMLString, "  ", " ")

	err = ioutil.WriteFile(configPath, []byte(updatedYAMLString), 0644)
	if err != nil {
		return err
	}

	return nil
}
