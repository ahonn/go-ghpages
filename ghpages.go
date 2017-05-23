package ghpages

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const CACHE_DIR = ".ghpages_cache"

func GetCacheDir() string {
	usr, err := user.Current()
	CheckIfErr(err)
	cacheDir := path.Join(usr.HomeDir, CACHE_DIR)
	return cacheDir
}

func CleanCacheDir() error {
	cacheDir := GetCacheDir()
	return os.RemoveAll(cacheDir)
}

func GetCloneDir(repo string) string {
	cacheDir := GetCacheDir()
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.Mkdir(cacheDir, 0777)
		CheckIfErr(err)
	}
	repoHash := base64.StdEncoding.EncodeToString([]byte(repo))
	cloneDir := path.Join(cacheDir, repoHash)
	return cloneDir
}

func GetAllAddFiles(dirPath string, matchFile string) []string {
	fi, err := os.Stat(dirPath)
	if os.IsNotExist(err) || !fi.IsDir() {
		log.Fatal("The base path option must be an existing directory")
	}

	// Get the files that need to be commited
	files := GetFilesList(dirPath)
	if len(files) == 0 {
		log.Fatal("The pattern in the 'src' property didn't match any files.")
	}
	return files
}

func Publish(basePath string, opt Options) {
	cacheDir := GetCacheDir()
	if opt.Clean {
		Info("Clean cache directory in: " + cacheDir)
		CleanCacheDir()
		return
	}

	repo := opt.GetRepo()
	cloneDir := GetCloneDir(repo)
	destDir := path.Join(cloneDir, opt.Dest)
	addFiles := GetAllAddFiles(basePath, opt.Src)

	git := &GitClient{
		Dir: cloneDir,
		Opt: opt,
	}

	// Clone repository
	Info("Clone " + repo)
	err := git.Clone(repo, cloneDir)
	if err != nil {
		Prompt(err)
	}

	// Remove untracked files form the working tree
	Info("Clean untracked files")
	err = git.Clean()
	CheckIfErr(err)

	// Download objects and refs from another repository
	Info("Fetch objects and refs")
	err = git.Fetch(opt.Remote)
	CheckIfErr(err)

	// Checkout the Branch
	Info("Checkout the branch to " + opt.Branch)
	err = git.Checkout(opt.Branch)
	CheckIfErr(err)

	// Remove files
	removeFiles, err := filepath.Glob(path.Join(destDir, "*"))
	if !opt.Add {
		Info("Remove destination files")
		CheckIfErr(err)
		if len(removeFiles) != 0 {
			err = git.Rm(removeFiles)
			CheckIfErr(err)
		}
	}
	err = git.RmFiles(removeFiles)
	CheckIfErr(err)

	// Copy files
	Info("Copy files into cache directory")
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err := os.Mkdir(destDir, 0777)
		CheckIfErr(err)
	}
	for _, file := range addFiles {
		rel, _ := filepath.Rel(basePath, file)
		dest := path.Join(destDir, rel)
		err = CopyFile(file, dest, true)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Add all files
	Info("Add all files")
	err = git.Add()
	CheckIfErr(err)

	// // Commit change with message
	Info("Commit change with message")
	err = git.Commit(opt.Message)
	CheckIfErr(err)

	// // Push a branch
	Info("Push to " + opt.Branch + " branch")
	err = git.Push(opt.Remote, opt.Branch)
	CheckIfErr(err)
}
