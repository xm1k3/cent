package jobs

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/xm1k3/cent/internal/utils"
)

func Start(_path string, name string, keepfolders bool, console bool) {
	timestamp := time.Now().Unix()

	for index, gitPath := range viper.GetStringSlice("community-templates") {
		utils.RunCommand("git clone "+gitPath+" /tmp/cent"+strconv.Itoa(int(timestamp))+"/repo"+strconv.Itoa(index), console)
	}

	dirname := "/tmp/cent" + strconv.Itoa(int(timestamp)) + "/"

	filepath.Walk(dirname,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			directory := getDirPath(strings.TrimPrefix(path, dirname))
			if info.IsDir() {
				for _, exDirs := range viper.GetStringSlice("exclude-dirs") {
					if strings.Contains(path, exDirs) {
						err := os.RemoveAll(path)
						if err != nil {
							log.Fatal(err)
						}
						return filepath.SkipDir
					} else {

						if keepfolders {
							if _, err := os.Stat(_path + "/" + name + directory); os.IsNotExist(err) {
								os.Mkdir(_path+"/"+name+directory, 0700)
							}
						}
						break
					}
				}
			} else {
				for _, exFiles := range viper.GetStringSlice("exclude-files") {
					if strings.Contains(path, exFiles) {
						e := os.Remove(path)
						if e != nil {
							log.Fatal(e)
						}
					} else {
						basename := info.Name()
						if filepath.Ext(basename) == ".yaml" {
							if !keepfolders {
								directory = ""
							}
							utils.RunCommand("cp "+path+" "+_path+"/"+name+directory, console)
						}
					}
					break
				}
			}
			return nil
		})

	DeleteFromTmp(dirname)
}

func UpdateRepo(path string, remDirs bool, remFiles bool) {
	filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if remDirs {
					for _, exDirs := range viper.GetStringSlice("exclude-dirs") {
						if strings.Contains(path, exDirs) {
							err := os.RemoveAll(path)
							if err != nil {
								log.Fatal(err)
							}
							fmt.Println(color.RedString("[D][-] Dir  removed\t" + path))
							return filepath.SkipDir
						}
					}
				}
			} else {
				if remFiles {
					for _, exFiles := range viper.GetStringSlice("exclude-files") {
						// fmt.Println("Path: ", path, exFiles)
						if strings.Contains(path, exFiles) {
							e := os.Remove(path)
							if e != nil {
								log.Fatal(e)
							}
							fmt.Println(color.RedString("[F][-] File removed\t" + path))
						}
						// break
					}
				}
			}
			return nil
		})
}

func getDirPath(path string) string {
	reponame := strings.Split(path, "/")[0]
	endpoint := strings.TrimPrefix(path, reponame)
	return endpoint
}

func RemoveEmptyFolders(dirname string) {

	f, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			if IsEmpty(dirname + "/" + file.Name()) {
				err := os.RemoveAll(dirname + "/" + file.Name())
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func IsEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false // Either not empty or error, suits both cases
}

func DeleteFromTmp(dirname string) {
	err := os.RemoveAll(dirname)
	if err != nil {
		log.Fatal(err)
	}
}
