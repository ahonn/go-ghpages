package ghpages

import (
  "log"
  "os"
  "os/exec"
  "path"
  "path/filepath"
  "strings"
  // "gopkg.in/src-d/go-git.v4"
)

type Config struct {
  Dist string
  Src  string
  Branch string
  Dest string
  Add bool
  Silent bool
  Message string
  Dotfiles bool
  Repo string
  Depth int
  Remote string
}

func getCachDir() string {
  pwd, err := os.Getwd()
  if err != nil {
    log.Fatal(err)
  }
  return path.Join(pwd, ".cache/")
}

func getRepo(config Config) string {
  if config.Repo != "" {
    return config.Repo
  }
  key := "remote." + config.Remote + ".url"
  out, err := exec.Command("git", "config", "--get", key).Output()
  if err != nil {
    log.Fatal(err)
  }
  output := string(out)
  return strings.Replace(output, "\n","",-1)
}

func Publish(basePath string, config Config) {
  log.Println(basePath)
  log.Println(config)

  // Exit when file is not exist or not directory.
  fi, err := os.Stat(basePath)
  if os.IsNotExist(err) || !fi.IsDir() {
    log.Fatal("The base path option must be an existing directory")
    os.Exit(1)
  }

  files, err := filepath.Glob(path.Join(basePath, config.Src))
  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }
  if len(files) == 0 {
    log.Fatal("The pattern in the 'src' property didn't match any files.")
    os.Exit(1)
  }

  repo := getRepo(config)
  cloneDir := getCachDir()

  log.Printf("clone %s to %s", repo, cloneDir)

  // log.Println(files)
}

