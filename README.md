# Go gh-pages
Tool for publishing files to a gh-pages branch on GitHub.
package inspiration [gh-pages](https://github.com/tschaub/gh-pages) and the options too.

## Usage

### as CLI
```
ghpages -d public -r git@github.com:ahonn/go-ghpages.git
```

### as a package
``` go
import "github.com/ahonn/go-ghpages"

opt := ghpages.Options {
  Dist:    "public",
  Branch:  "gh-pages",
  Dest:    ".",
  Add:     false,
  Message: "update messages",
  Repo:    "git@github.com:ahonn/go-ghpages.git",
  Depth:   1,
  Remote:  "origin",
}

ghpages.Public("public", opt)
```

## LICENSE
see [MIT LICENSE](./LICENSE)
