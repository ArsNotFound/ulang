package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/Ars2014/ulang/ast"
	"github.com/Ars2014/ulang/eval"
	"github.com/Ars2014/ulang/lexer"
	"github.com/Ars2014/ulang/object"
	"github.com/Ars2014/ulang/parser"
)

const Prompt = ">>> "

const ULang = `   __  ____    ___    _   ________
  / / / / /   /   |  / | / / ____/
 / / / / /   / /| | /  |/ / / __  
/ /_/ / /___/ ___ |/ /|  / /_/ /  
\____/_____/_/  |_/_/ |_/\____/
`

type Options struct {
	Debug       bool
	Interactive bool
}

type REPL struct {
	user string
	args []string
	opts *Options
}

func New(user string, args []string, opts *Options) *REPL {
	object.Stdin = os.Stdin
	object.Stdout = os.Stdout
	object.ExitFn = os.Exit

	return &REPL{user, args, opts}
}

func (r *REPL) Eval(f io.Reader) (env *object.Environment) {
	env = object.NewEnvironment()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading source file: %s", err)
		return
	}

	l := lexer.NewLexer(b)
	p := parser.NewParser()

	program, err := p.Parse(l)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error occured while parsing program: %s", err)
		return
	}

	eval.Eval(program.(*ast.Program), env)
	return
}

func (r *REPL) StartEvalLoop(in io.Reader, out io.Writer, env *object.Environment) {
	scanner := bufio.NewScanner(in)

	if env == nil {
		env = object.NewEnvironment()
	}

	p := parser.NewParser()

	for {
		fmt.Printf(Prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Bytes()
		line = bytes.Trim(line, " \t\n\r")

		if len(line) == 0 {
			continue
		}

		l := lexer.NewLexer(line)
		p.Reset()

		program, err := p.Parse(l)
		if err != nil {
			fmt.Printf("error occured while parsing program: %s\n", err)
			continue
		}

		obj := eval.Eval(program.(*ast.Program), env)
		if obj != nil {
			if _, ok := obj.(*object.Null); !ok {
				io.WriteString(out, obj.Inspect()+"\n")
			}
		}
	}
}

func (r *REPL) Run() {
	object.Arguments = make([]string, len(r.args))
	copy(object.Arguments, r.args)

	if len(r.args) == 0 {
		fmt.Printf(ULang)
		fmt.Printf("Hello %s! This is the ULang programming language!\n", r.user)
		fmt.Printf("Feel free to type in commands\n")
		r.StartEvalLoop(os.Stdin, os.Stdout, nil)
		return
	}

	if len(r.args) > 0 {
		f, err := os.Open(r.args[0])
		if err != nil {
			log.Fatalf("could not open source file %s: %s", r.args[0], err)
		}

		r.args = r.args[1:]
		object.Arguments = object.Arguments[1:]

		env := r.Eval(f)
		if r.opts.Interactive {
			r.StartEvalLoop(os.Stdin, os.Stdout, env)
		}
	}
}
