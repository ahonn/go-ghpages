package ghpages

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

const CACHE_DIR = ".ghpages_cache"

type Options struct {
	Dist     string
	Src      string
	Branch   string
	Dest     string
	Add      bool
	Silent   bool
	Message  string
	Dotfiles bool
	Repo     string
	Depth    string
	Remote   string
	Clean    bool
}

// Get repository url. if unset `Repo` filed, exec `git config`
func (this *Options) GetRepo() string {
	if this.Repo != "" {
		return this.Repo
	}
	args := []string{
		"config",
		"--get",
		"remote." + this.Remote + ".url",
	}
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	output := strings.Replace(string(out), "\n", "", -1)
	return output
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getCacheDir() string {
	usr, err := user.Current()
	handleError(err)
	cacheDir := path.Join(usr.HomeDir, CACHE_DIR)
	return cacheDir
}

func cleanCacheDir() error {
	cacheDir := getCacheDir()
	return os.RemoveAll(cacheDir)
}

func getCloneDir(repo string) string {
	cacheDir := getCacheDir()
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.Mkdir(cacheDir, 0777)
		handleError(err)
	}
	repoHash := base64.StdEncoding.EncodeToString([]byte(repo))
	cloneDir := path.Join(cacheDir, repoHash)
	return cloneDir
}

func Publish(basePath string, opt Options) {
	fmt.Println(basePath)
	fmt.Println(opt)

	cacheDir := getCacheDir()
	if opt.Clean {
		fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Clean cache directory in: "+cacheDir)
		cleanCacheDir()
		return
	}

	// Exit when file is not exist or not directory.
	fi, err := os.Stat(basePath)
	if os.IsNotExist(err) || !fi.IsDir() {
		log.Fatal("The base path option must be an existing directory")
	}

	// Get the files that need to be commited
	// TODO: Fix this, files is not all files
	files, err := filepath.Glob(path.Join(basePath, opt.Src))
	handleError(err)
	if len(files) == 0 {
		log.Fatal("The pattern in the 'src' property didn't match any files.")
	}

	repo := opt.GetRepo()
	cloneDir := getCloneDir(repo)
	destDir := path.Join(cloneDir, opt.Dest)

	git := &GitClient{
		Dir: cloneDir,
		Opt: opt,
	}

	// Check the remote if clone directory is exists
	if _, err := os.Stat(cloneDir); !os.IsNotExist(err) {
		check, err := git.CheckRemote()
		handleError(err)
		if !check {
			logMsg := "Clean repository cache directory: cache remote is not " + repo
			fmt.Printf("\x1b[34;1m%s\x1b[0m\n", logMsg)
			os.RemoveAll(cloneDir)
		}
	}

	// Clone repository
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Clone "+repo)
	err = git.Clone(repo, cloneDir)
	if err != nil {
		fmt.Printf("\x1b[36;1m%s\x1b[0m\n", err)
	}

	// Remove untracked files form the working tree
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Clean untracked files")
	err = git.Clean()
	handleError(err)

	// Download objects and refs from another repository
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Fetch objects and refs")
	err = git.Fetch(opt.Remote)
	handleError(err)

	// Checkout the Branch
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Checkout the branch to "+opt.Branch)
	err = git.Checkout(opt.Branch)
	handleError(err)

	// Remove files
	if !opt.Add {
		fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Remove destination files")
		removeFiles, err := filepath.Glob(path.Join(destDir, "*"))
		handleError(err)
		if len(removeFiles) != 0 {
			err = git.Rm(removeFiles)
			handleError(err)
		}
	}

	// Copy files
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Copy files into cache directory")
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err := os.Mkdir(destDir, 0777)
		handleError(err)
	}
	for _, file := range files {
		fmt.Println(file)
		base := filepath.Base(file)
		newName := path.Join(cloneDir, opt.Dest, base)
		err := os.Link(file, newName)
		if err != nil {
			pwd, _ := os.Getwd()
			rel, _ := filepath.Rel(pwd, file)
			fmt.Printf("\x1b[36;1m%s\x1b[0m\n", "Skip copy file: "+rel)
		}
	}

	// Add all files
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Add all files")
	err = git.Add()
	handleError(err)

	// Commit change with message
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Commit change with message")
	err = git.Commit(opt.Message)
	handleError(err)

	// Push a branch
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Push to "+opt.Branch+" branch")
	err = git.Push(opt.Remote, opt.Branch)
	handleError(err)
}
