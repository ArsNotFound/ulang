package object

import "io"

var (
	Arguments []string
	Stdin     io.Reader
	Stdout    io.Writer
	ExitFn    func(int)
)
