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

	// TODO: allow flag to override config file?

	cfg, err := zengarden.LoadConfig(zengarden.DefaultConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while reading config: %v\n", err)
	}

	if err := zengarden.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error while building: %v\n", err)
	}
}
