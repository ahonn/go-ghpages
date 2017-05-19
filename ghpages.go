package ghpages

import (
  "log"
  "os"
  "os/exec"
  "path"
  "path/filepath"
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

func getRepo(config Config) string {
  if config.Repo != "" {
    return config.Repo
  }
  cmd := "git config --get remote." + config.Remote + ".url"
  out, err := exec.Command(cmd).Output()
  if err != nil {
    log.Fatal(err)
  }
  return string(out)
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
  log.Println(repo)
  log.Println(files)
}

