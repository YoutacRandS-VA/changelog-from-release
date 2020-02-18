package main

import (
	"flag"
	"fmt"
	"os"
)

var version = "v1.2.0"

func usage() {
	fmt.Fprint(os.Stderr, "Usage: changelog-from-release [flags]\n\n")
	flag.PrintDefaults()
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	os.Exit(111)
}

func main() {
	flag.Usage = usage
	ver := flag.Bool("v", false, "Output version to stdout")
	tag := flag.Bool("t", false, "Output the latest tag value to stdout after generating changelog")
	flag.Parse()

	if *ver {
		fmt.Println(version)
		os.Exit(0)
	}

	if flag.NArg() != 0 {
		usage()
		os.Exit(111)
	}

	git, err := NewGitForCwd()
	if err != nil {
		fail(err)
	}

	url, err := git.TrackingRemoteURL()
	if err != nil {
		fail(err)
	}

	gh, err := GitHubFromURL(url)
	if err != nil {
		fail(err)
	}

	rels, err := gh.Releases()
	if err != nil {
		fail(err)
	}
	if len(rels) == 0 {
		fail(fmt.Errorf("No release was found at %s", url))
	}

	cl, err := NewChangeLog(git.root, url)
	if err != nil {
		fail(err)
	}

	if err := cl.Generate(rels); err != nil {
		fail(err)
	}

	if *tag {
		fmt.Println(rels[0].GetTagName())
	}
}
