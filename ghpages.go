package ghpages

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const CACHE_DIR = ".ghpages_cache"

type Options struct {
	Dist     string
	Src      string
	Branch   string
	Dest     string
	Silent   bool
	Message  string
	Dotfiles bool
	Repo     string
	Depth    string
	Remote   string
	Clean    bool
}

func getCacheDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cacheDir := path.Join(pwd, usr.HomeDir, CACHE_DIR)
	return cacheDir
}

func cleanCacheDir() error {
	cacheDir := getCacheDir()
	return os.RemoveAll(cacheDir)
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

	files, err := filepath.Glob(path.Join(basePath, opt.Src))
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		log.Fatal("The pattern in the 'src' property didn't match any files.")
	}

	git := &GitClient{
		Dir: cacheDir,
		Opt: opt,
	}
	repo := git.GetRepo()

	if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
		check, err := git.CheckRemote()
		if err != nil {
			log.Fatal(err)
		}
		if !check {
			fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Clean cache directory: cache remote is not "+repo)
			cleanCacheDir()
		}
	}

	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Clone "+repo+" to "+cacheDir)
	err = git.Clone(repo, cacheDir)
	if err != nil {
		fmt.Printf("\x1b[36;1m%s\x1b[0m\n", err)
	}
}
