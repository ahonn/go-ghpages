package main

import (
  "os"
  "path"
  "log"
  "github.com/urfave/cli"
  "github.com/ahonn/go-ghpages"
)

func main() {
  app := cli.NewApp()

  app.Flags = []cli.Flag {
    cli.StringFlag {
      Name: "dist, d",
      Value: ".",
      Usage: "Base durectory for all source `files`",
    },
    cli.StringFlag {
      Name: "src, s",
      Value: "**/*",
      Usage: "Pattern used to select which `files` to publish",
    },
    cli.StringFlag {
      Name: "branch, b",
      Value: "gh-pages",
      Usage: "Name of the `branch` you are pushing to",
    },
    cli.StringFlag {
      Name: "dest, e",
      Value: ".",
      Usage: "Target `directory` within the destination branch (relative to the root)",
    },
    cli.BoolFlag {
      Name: "add, a",
      Usage: "Only add, and never remove existing files",
    },
    cli.BoolFlag {
      Name: "silent, x",
      Usage: "Do not output the repository url",
    },
    cli.StringFlag {
      Name: "message, m",
      Value: "Updates",
      Usage: "Commit `message`",
    },
    cli.BoolFlag {
      Name: "dotfiles, t",
      Usage: "Include dotfiles",
    },
    cli.StringFlag {
      Name: "repo, r",
      Usage: "URL of the `repository` you are pushing to",
    },
    cli.IntFlag {
      Name: "depth, p",
      Value: 1,
      Usage: "Depth for clone",
    },
    cli.StringFlag {
      Name: "remote, o",
      Value: "origin",
      Usage: "The `name` of the remote",
    },
  }

  app.Action = func(c *cli.Context) error {
    pwd, err := os.Getwd()
    if err != nil {
      log.Fatal(err)
    }

    ghpages.Publish(path.Join(pwd, c.String("dist")), ghpages.Config {
      Dist: c.String("dist"),
      Src: c.String("src"),
      Branch: c.String("branch"),
      Dest: c.String("dest"),
      Add: c.Bool("add"),
      Silent: c.Bool("silent"),
      Message: c.String("message"),
      Dotfiles: c.Bool("dotfiles"),
      Repo: c.String("repo"),
      Depth: c.Int("depth"),
      Remote: c.String("remote"),
    })

    return nil
  }

  app.Run(os.Args)
}
