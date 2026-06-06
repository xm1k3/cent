package jobs

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	filepath "path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"github.com/xm1k3/cent/v2/internal/utils"
)

type RepoEntry struct {
	URL     string   `yaml:"url"`
	Commit  string   `yaml:"commit,omitempty"`
	Exclude []string `yaml:"exclude,omitempty"`
}

func ParseRepoEntries() []RepoEntry {
	raw := viper.Get("community-templates")
	items, ok := raw.([]interface{})
	if !ok {
		return nil
	}

	var entries []RepoEntry
	for _, item := range items {
		switch v := item.(type) {
		case string:
			entries = append(entries, RepoEntry{URL: v})
		case map[string]interface{}:
			entry := RepoEntry{}
			if u, ok := v["url"].(string); ok {
				entry.URL = u
			}
			if c, ok := v["commit"].(string); ok {
				entry.Commit = c
			}
			if e, ok := v["exclude"].([]interface{}); ok {
				for _, ex := range e {
					if s, ok := ex.(string); ok {
						entry.Exclude = append(entry.Exclude, s)
					}
				}
			}
			entries = append(entries, entry)
		}
	}
	return entries
}

func cloneRepo(entry RepoEntry, console bool, index string, timestamp string, timeoutMinutes int) error {
	destDir := filepath.Join(os.TempDir(), fmt.Sprintf("cent%s/repo%s", timestamp, index))

	if err := os.MkdirAll(destDir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	timeout := time.Duration(timeoutMinutes) * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cloneArgs := []string{"clone", "--single-branch", "--no-tags", "--no-recurse-submodules"}
	if entry.Commit == "" {
		cloneArgs = append(cloneArgs, "--depth", "1")
	}
	cloneArgs = append(cloneArgs, entry.URL, destDir)

	cmd := exec.CommandContext(ctx, "git", cloneArgs...)

	if console {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("git clone timed out for %s", entry.URL)
	}
	if err != nil {
		return fmt.Errorf("git clone failed for %s: %w", entry.URL, err)
	}

	if entry.Commit != "" {
		cmdCheckout := exec.CommandContext(ctx, "git", "-C", destDir, "checkout", entry.Commit)
		if console {
			cmdCheckout.Stdout = os.Stdout
			cmdCheckout.Stderr = os.Stderr
		}
		if err := cmdCheckout.Run(); err != nil {
			return fmt.Errorf("git checkout %s failed for %s: %w", entry.Commit, entry.URL, err)
		}
		fmt.Printf(color.GreenString("[CLONED] %s @ %s\n", entry.URL, entry.Commit))
	} else {
		fmt.Printf(color.GreenString("[CLONED] %s\n", entry.URL))
	}
	return nil
}

type repoWork struct {
	index string
	entry RepoEntry
}

func worker(work chan repoWork, wg *sync.WaitGroup, console bool, timestamp string, timeoutMinutes int) {
	defer wg.Done()
	for repo := range work {
		err := cloneRepo(repo.entry, console, repo.index, timestamp, timeoutMinutes)
		if err != nil {
			fmt.Println(color.RedString("[ERR] clone: " + repo.entry.URL + " - " + err.Error()))
		}
	}
}

func repoNameFromURL(url string) string {
	url = strings.TrimSuffix(url, ".git")
	url = strings.TrimSuffix(url, "/")
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return url
}

func Start(_path string, console bool, threads int, defaultTimeout int, keepFolders bool, byRepo bool) {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	if _, err := os.Stat(filepath.Join(_path)); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(_path), 0700)
	}

	if _, err := exec.LookPath("git"); err != nil {
		log.Fatalf("Git is not installed or not available in PATH: %v", err)
	}

	entries := ParseRepoEntries()
	repoNames := make(map[string]string)
	repoExcludes := make(map[string][]string)
	for index, entry := range entries {
		idx := strconv.Itoa(index)
		repoNames[idx] = repoNameFromURL(entry.URL)
		if len(entry.Exclude) > 0 {
			repoExcludes[idx] = entry.Exclude
		}
	}

	work := make(chan repoWork)
	go func() {
		for index, entry := range entries {
			work <- repoWork{index: strconv.Itoa(index), entry: entry}
		}
		close(work)
	}()

	wg := &sync.WaitGroup{}

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(work, wg, console, timestamp, defaultTimeout)
	}
	wg.Wait()

	dirname := filepath.Join(os.TempDir(), "cent"+timestamp)

	filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath := strings.TrimPrefix(path, dirname+string(os.PathSeparator))
		parts := strings.SplitN(relPath, string(os.PathSeparator), 2)
		repoDir := parts[0]
		idx := strings.TrimPrefix(repoDir, "repo")

		if excludes, ok := repoExcludes[idx]; ok {
			for _, ex := range excludes {
				if strings.Contains(relPath, ex) {
					return nil
				}
			}
		}

		if filepath.Ext(info.Name()) != ".yaml" {
			return nil
		}

		var destinationPath string
		if keepFolders {
			innerPath := ""
			if len(parts) > 1 {
				innerPath = parts[1]
			}
			if !strings.Contains(innerPath, string(os.PathSeparator)) {
				innerPath = filepath.Join("others", innerPath)
			}
			if byRepo {
				repoName := repoNames[idx]
				innerPath = filepath.Join(repoName, innerPath)
			}
			destinationPath = filepath.Join(_path, innerPath)
			os.MkdirAll(filepath.Dir(destinationPath), 0700)
		} else {
			destinationPath = filepath.Join(_path, info.Name())
		}

		if cpErr := utils.CopyFile(path, destinationPath); cpErr != nil {
			return cpErr
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

	return err == io.EOF // Either not empty or error, suits both cases
}

func DeleteFromTmp(dirname string) {
	err := os.RemoveAll(dirname)
	if err != nil {
		log.Fatal(err)
	}
}
