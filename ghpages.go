package ghpages

import (
  "log"
  "os"
  "os/exec"
  "os/user"
  "path"
  "path/filepath"
  "strings"
  "io/ioutil"

  "golang.org/x/crypto/ssh"
  "gopkg.in/src-d/go-git.v4"
  gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
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

func getCacheDir() string {
  pwd, err := os.Getwd()
  if err != nil {
    log.Fatal(err)
  }
  return path.Join(pwd, ".cache/")
}

func getRepo(c Config) string {
  if c.Repo != "" {
    return c.Repo
  }
  key := "remote." + c.Remote + ".url"
  out, err := exec.Command("git", "config", "--get", key).Output()
  if err != nil {
    log.Fatal(err)
  }
  output := strings.Replace(string(out), "\n", "", -1)
  return output
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

  cloneDir := getCacheDir()
  repo := getRepo(config)

  usr, err := user.Current()
  if err != nil {
      log.Fatal( err )
  }
  key, err := ioutil.ReadFile(path.Join(usr.HomeDir, ".ssh/id_rsa"))
  if err != nil {
    log.Fatalf("unable to read private key: %v", err)
  }
  signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
  }

  log.Printf("\x1b[34;1m%s\x1b[0m\n", "Clone " + repo + " to " + cloneDir)
  _, err = git.PlainClone(cloneDir, false, &git.CloneOptions {
    URL: repo,
    Depth: config.Depth,
    SingleBranch: true,
    Auth: &gitssh.PublicKeys{User: "git", Signer: signer},
    RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
    Progress: os.Stdout,
  })
  if err != nil {
    log.Fatal(err)
  }

  // log.Println(files)
}

