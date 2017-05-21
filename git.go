package ghpages

import (
  "os"
  "os/exec"
  "path/filepath"
  "log"
  "errors"
)

type GitClient struct {
  Dir string
  Opt Options
}

func (g *GitClient) Command(dir string, args ...string) *exec.Cmd {
  name := "git"

  if dir == "" {
    dir = g.Dir
  }

  cmd := &exec.Cmd {
    Path: name,
    Args: append([]string{name}, args...),
    Dir: dir,
  }

  if filepath.Base(name) == name {
    if lp, err := exec.LookPath(name); err != nil {
      log.Fatal(err)
    } else {
      cmd.Path = lp
    }
  }
  return cmd
}

func (g *GitClient) Clone(repo string, cloneDir string) error {
  if _, err := os.Stat(cloneDir); os.IsNotExist(err) {
    opt := g.Opt
    args := []string {
      "clone",
      repo,
      cloneDir,
      "--branch",
      opt.Branch,
      "--single-branch",
      "--origin",
      opt.Remote,
      "--depth",
      opt.Depth,
    }
    cmd := g.Command(".", args...)
    return cmd.Run()
  }
  return errors.New("Repository is already exists")
}

