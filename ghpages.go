package ghpages

import (
  "log"
  "fmt"
  "os"
  "os/exec"
  "os/user"
  "path"
  "path/filepath"
  "strings"
  "io/ioutil"

  "golang.org/x/crypto/ssh"
  "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing"
  gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type Options struct {
  Dist string
  Src  string
  Branch string
  Dest string
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

func getRepo(opt Options) string {
  if opt.Repo != "" {
    return opt.Repo
  }
  key := "remote." + opt.Remote + ".url"
  out, err := exec.Command("git", "config", "--get", key).Output()
  if err != nil {
    log.Fatal(err)
  }
  output := strings.Replace(string(out), "\n", "", -1)
  return output
}

func getSSHSigner() ssh.Signer {
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
  return signer
}

func Publish(basePath string, opt Options) {
  fmt.Println(basePath)
  fmt.Println(opt)

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

  cloneDir := getCacheDir()
  repo := getRepo(opt)
  signer := getSSHSigner()

  branch := plumbing.ReferenceName("refs/heads/" + opt.Branch)

  // Clone Repo
  fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Clone " + repo + " to " + cloneDir)
  r, err := git.PlainClone(cloneDir, false, &git.CloneOptions {
    URL: repo,
    Depth: opt.Depth,
    RemoteName: opt.Remote,
    ReferenceName: branch,
    Auth: &gitssh.PublicKeys{User: "git", Signer: signer},
    Progress: os.Stdout,
  })
  if err != nil {
    fmt.Printf("\x1b[36;1m%s\x1b[0m\n", err)
    r, err = git.PlainOpen(cloneDir)
    if err != nil {
      log.Fatal(err)
    }
  }

  // Checkout Repo
  fmt.Printf("\x1b[34;1m%s\x1b[0m\n", "Checkout to branch: " + opt.Branch)
  w, err := r.Worktree()
  if err != nil {
    log.Fatal(err)
  }
  err = w.Checkout(&git.CheckoutOptions {
    Branch: plumbing.ReferenceName("refs/heads/" + opt.Branch),
  })
  if err != nil {
    log.Fatal(err)
  }

  // log.Println(files)
}

