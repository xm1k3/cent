package jobs

import (
	"fmt"
	"io"
	"log"
	"os"
	filepath "path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/xm1k3/cent/internal/utils"
)

func Start(_path string, keepfolders bool, console bool) {
	timestamp := time.Now().Unix()
	if _, err := os.Stat(filepath.Join(_path)); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(_path), 0700)
	}

	for index, gitPath := range viper.GetStringSlice("community-templates") {
		utils.RunCommand("git clone "+gitPath+" /tmp/cent"+strconv.Itoa(int(timestamp))+"/repo"+strconv.Itoa(index), console)
		if !console {
			fmt.Println(color.GreenString("[CLONED] \t" + gitPath))
		}
	}

	dirname := "/tmp/cent" + strconv.Itoa(int(timestamp)) + "/"

	filepath.Walk(dirname,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			directory := getDirPath(strings.TrimPrefix(path, dirname))
			if info.IsDir() {
				if keepfolders {
					if _, err := os.Stat(filepath.Join(_path, directory)); os.IsNotExist(err) {
						// fmt.Println(_path, name, directory)
						os.Mkdir(filepath.Join(_path, directory), 0700)
					}
				}
			} else {
				basename := info.Name()
				if filepath.Ext(basename) == ".yaml" {
					if !keepfolders {
						directory = ""
					}
					utils.RunCommand("cp "+path+" "+filepath.Join(_path, directory), console)
				}
			}
			return nil
		})

	DeleteFromTmp(dirname)
}

func UpdateRepo(path string, remDirs bool, remFiles bool, printOut bool) {
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
							if printOut {
								fmt.Println(color.RedString("[D][-] Dir  removed\t" + path))
							}
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
							if printOut {
								fmt.Println(color.RedString("[F][-] File removed\t" + path))
							}

						}
						// break
					}
				}
			}
			return nil
		})
}

func RemoveDuplicates(path string, console bool) {
	utils.RunCommand("fdupes -d -N -r "+path, console)
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
			if IsEmpty(filepath.Join(dirname, file.Name())) {
				err := os.RemoveAll(filepath.Join(dirname, file.Name()))
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
