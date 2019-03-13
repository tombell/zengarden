package main

import (
	"flag"
	"fmt"
	"log"
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

func validateFlags() {
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "zen-garden %s (%s)\n", Version, Commit)
		os.Exit(0)
	}

	logger := log.New(os.Stderr, "[zen] ", log.LstdFlags)

	if err := zengarden.Run(); err != nil {
		logger.Fatalf("error while running: %v\n", err)
	}
}
