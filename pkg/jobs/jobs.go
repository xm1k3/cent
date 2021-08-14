package jobs

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	filepath "path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/xm1k3/cent/internal/utils"
)

func cloneRepo(gitPath string, console bool, index string, timestamp string) {
	utils.RunCommand("git clone "+gitPath+" /tmp/cent"+timestamp+"/repo"+index, console)
	if !console {
		fmt.Println(color.GreenString("[CLONED] \t" + gitPath))
	}
}

func worker(work chan [2]string, wg *sync.WaitGroup, console bool, timestamp string) {
	defer wg.Done()
	for repo := range work {
		cloneRepo(repo[1], console, repo[0], timestamp)
	}
}

func Start(_path string, keepfolders bool, console bool, threads int) {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	if _, err := os.Stat(filepath.Join(_path)); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(_path), 0700)
	}

	work := make(chan [2]string)
	go func() {
		for index, gitPath := range viper.GetStringSlice("community-templates") {
			work <- [2]string{strconv.Itoa(index), gitPath}
		}
		close(work)
	}()

	wg := &sync.WaitGroup{}

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(work, wg, console, timestamp)
	}
	wg.Wait()

	dirname := "/tmp/cent" + timestamp + "/"

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
	fmt.Println("Removing duplicate templates...")
	files := getFilePaths(path)
	hashes := make(map[string]string)

	// get file hashes
	for _, file := range files {
		hashes[file] = getFileHash(file)
	}

	// get hash files
	hashfiles := make(map[string][]string)
	for file, hash := range hashes {
		hashfiles[hash] = append(hashfiles[hash], file)
	}

	// for each hash, remove all the files except the first one
	for _, files := range hashfiles {
		sort.Strings(files)
		for _, fileToRemove := range files[1:] {
			if console {
				fmt.Printf("Removing duplicate file: %s\n", fileToRemove)
			}
			os.Remove(fileToRemove)
		}
	}
}

func getFilePaths(path string) []string {
	var files []string

	// go through each file
	err := filepath.WalkDir(path, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !d.IsDir() {
			files = append(files, s)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func getFileHash(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
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
