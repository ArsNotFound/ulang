package main

import "fmt"

var (
	Version = "0.0.1"

	GitCommit = "HEAD"
)

func FullVersion() string {
	return fmt.Sprintf("%s@%s", Version, GitCommit)
}
