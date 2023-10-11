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
	"path"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xm1k3/cent/pkg/jobs"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "cent",
	Short: "Community edition nuclei templates",
	Long: `
 ██████╗███████╗███╗   ██╗████████╗
██╔════╝██╔════╝████╗  ██║╚══██╔══╝
██║     █████╗  ██╔██╗ ██║   ██║   
██║     ██╔══╝  ██║╚██╗██║   ██║   
╚██████╗███████╗██║ ╚████║   ██║   
 ╚═════╝╚══════╝╚═╝  ╚═══╝   ╚═╝   
									   	
Community edition nuclei templates, a simple tool that allows you 
to organize all the Nuclei templates offered by the community in one place.

Disclaimer: The developer of this tool is not responsible for how the community 
uses the open source nuclei-templates collected within it. 
These templates have not been validated by Project Discovery and are provided as-is.

By xm1k3`,
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		pathFlag, _ := cmd.Flags().GetString("path")
		console, _ := cmd.Flags().GetBool("console")
		threads, _ := cmd.Flags().GetInt("threads")
		timeout, _ := cmd.Flags().GetInt("timeout")

		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		_, errHome := os.Stat(path.Join(home, ".cent.yaml"))
		_, errDefault := os.Stat(cfgFile)
		if os.IsNotExist(errHome) && os.IsNotExist(errDefault) {
			fmt.Println(`Run ` + color.YellowString("cent init") + ` to automatically download ` +
				color.HiCyanString(".cent.yaml") + ` from repo and copy it to ` +
				color.HiCyanString("$HOME/.cent.yaml"))
			return
		}

		fmt.Println(color.CyanString("cent started"))

		jobs.Start(pathFlag, console, threads, timeout)
		jobs.RemoveEmptyFolders(path.Join(pathFlag))
		jobs.UpdateRepo(path.Join(pathFlag), true, true, false)
		jobs.RemoveDuplicates(path.Join(pathFlag), console)
		fmt.Println(color.YellowString("[!] Removed duplicates"))
		fmt.Println(color.CyanString("cent finished, you can find all your nuclei-templates in " + pathFlag))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cent.yaml)")

	rootCmd.Flags().StringP("path", "p", "cent-nuclei-templates", "Root path to save the templates")
	rootCmd.Flags().BoolP("console", "C", false, "Print console output")
	rootCmd.Flags().IntP("threads", "t", 10, "Number of threads to use when cloning repositories")
	rootCmd.Flags().IntP("timeout", "T", 2, "timeout in seconds")

	rootCmd.MarkFlagRequired("name")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".cent")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
