package ghpages

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GitClient struct {
	Dir string
	Opt Options
}

func (this *GitClient) Command(dir string, args ...string) *exec.Cmd {
	name := "git"
	if dir == "" {
		dir = this.Dir
	}

	cmd := &exec.Cmd{
		Path: name,
		Args: append([]string{name}, args...),
		Dir:  dir,
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

func (this *GitClient) Clean() error {
	args := []string{
		"clean",
		"-f",
		"-d",
	}
	cmd := this.Command("", args...)
	_, err := cmd.Output()
	return err
}

func (this *GitClient) Rm(removeFiles []string) error {
	files := []string{}
	for _, file := range removeFiles {
		if filepath.Base(file) != ".git" {
			files = append(files, file)
		}
	}

	args := []string{
		"rm",
		"--ignore-unmatch",
		"-r",
		"-f",
		strings.Join(files, " "),
	}
	cmd := this.Command("", args...)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return err
}

func (this *GitClient) CheckRemote() (bool, error) {
	repoRemote := this.Opt.GetRepo()

	args := []string{
		"config",
		"--get",
		"remote." + this.Opt.Remote + ".url",
	}
	cmd := this.Command("", args...)
	out, err := cmd.Output()
	cloneRemote := strings.Replace(string(out), "\n", "", -1)

	check := cloneRemote == repoRemote
	return check, err
}

func (this *GitClient) Fetch(remote string) error {
	args := []string{
		"fetch",
		remote,
	}
	cmd := this.Command("", args...)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (this *GitClient) Clone(repo string, cloneDir string) error {
	if _, err := os.Stat(cloneDir); os.IsNotExist(err) {
		args := []string{
			"clone",
			repo,
			cloneDir,
			"--branch",
			this.Opt.Branch,
			"--single-branch",
			"--origin",
			this.Opt.Remote,
			"--depth",
			this.Opt.Depth,
		}
		cmd := this.Command(".", args...)
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("\x1b[36;1m%s\x1b[0m\n", "Clone repository for master branch")
			args := []string{
				"clone",
				repo,
				cloneDir,
				"--origin",
				this.Opt.Remote,
			}
			cmd := this.Command(".", args...)
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			return err
		}
		return nil
	}
	return errors.New("Repository is already exists")
}

func (this *GitClient) Checkout(branch string) error {
	args := []string{
		"checkout",
		branch,
	}
	cmd := this.Command("", args...)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		args := []string{
			"checkout",
			"-b",
			branch,
		}
		cmd := this.Command("", args...)
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		return err
	}
	return nil
}

func (this *GitClient) Add() error {
	args := []string{
		"add",
		"--all",
	}
	cmd := this.Command("", args...)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func (this *GitClient) Commit(message string) error {
	args := []string{
		"commit",
		"-m",
		message,
	}
	cmd := this.Command("", args...)
	err := cmd.Run()
	return err
}

func (this *GitClient) Push(remote string, branch string) error {
	args := []string{
		"push",
		"--tags",
		remote,
		branch,
	}
	cmd := this.Command("", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return err
}
