package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/Ars2014/ulang/repl"
)

var (
	interactive bool
	version     bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [<filename>]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.BoolVar(&version, "v", false, "display version information")
	flag.BoolVar(&interactive, "i", false, "enable interactive mode")
}

func main() {
	flag.Parse()

	if version {
		fmt.Printf("%s %s", path.Base(os.Args[0]), FullVersion())
		os.Exit(0)
	}

	currUser, err := user.Current()
	if err != nil {
		log.Fatalf("could not determine current user: %s", err)
	}

	args := flag.Args()

	opts := &repl.Options{
		Debug:       false,
		Interactive: interactive,
	}
	repl_ := repl.New(currUser.Username, args, opts)
	repl_.Run()
}
