package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tombell/zengarden"
)

const helpText = `usage: zg [args]

Special options:
  --help      show this message, then exit
  --version   show the version number, then exit
`

var (
	version = flag.Bool("version", false, "")
)

func usage() {
	fmt.Fprintf(os.Stderr, helpText)
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "zen-garden %s (%s)\n", Version, Commit)
		os.Exit(0)
	}

	cfg := &zengarden.Config{
		Source:    ".",
		Target:    "_site",
		Permalink: "/post/:title/",
		Excludes:  []string{"README.md"},
	}

	if err := zengarden.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error while running: %v\n", err)
	}
}
