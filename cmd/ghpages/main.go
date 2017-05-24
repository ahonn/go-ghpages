package main

import (
	"github.com/ahonn/go-ghpages"
	"github.com/urfave/cli"
	"log"
	"os"
	"path"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dist, d",
			Value: ".",
			Usage: "base durectory for all source `files`",
		},
		cli.StringFlag{
			Name:  "branch, b",
			Value: "gh-pages",
			Usage: "name of the `branch` you are pushing to",
		},
		cli.StringFlag{
			Name:  "dest, e",
			Value: ".",
			Usage: "target `directory` within the destination branch (relative to the root)",
		},
		cli.BoolFlag{
			Name:  "add, a",
			Usage: "only add, and never remove existing files",
		},
		cli.StringFlag{
			Name:  "message, m",
			Value: "Updates",
			Usage: "commit `message`",
		},
		cli.StringFlag{
			Name:  "repo, r",
			Usage: "url of the `repository` you are pushing to",
		},
		cli.StringFlag{
			Name:  "depth, p",
			Value: "1",
			Usage: "depth for clone",
		},
		cli.StringFlag{
			Name:  "remote, o",
			Value: "origin",
			Usage: "the `name` of the remote",
		},
		cli.BoolFlag{
			Name:  "clean, c",
			Usage: "clean cache directory",
		},
	}

	app.Action = func(c *cli.Context) error {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		ghpages.Publish(path.Join(pwd, c.String("dist")), ghpages.Options{
			Dist:    c.String("dist"),
			Branch:  c.String("branch"),
			Dest:    c.String("dest"),
			Add:     c.Bool("add"),
			Message: c.String("message"),
			Repo:    c.String("repo"),
			Depth:   c.String("depth"),
			Remote:  c.String("remote"),
			Clean:   c.Bool("clean"),
		})

		return nil
	}

	app.Run(os.Args)
}
