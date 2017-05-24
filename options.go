package ghpages

import (
	"log"
	"os/exec"
	"strings"
)

type Options struct {
	Dist    string
	Branch  string
	Dest    string
	Add     bool
	Message string
	Repo    string
	Depth   string
	Remote  string
	Clean   bool
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
