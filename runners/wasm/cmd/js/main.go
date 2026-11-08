// Command js is a JavaScript engine (goja, ES2015+) compiled to WASI with a
// small Node-like global set: console, process, and require("fs").readFileSync
// for stdin. The grading host runs it as `js /work/<entry>`.
//
// Exit codes: process.exit's code, 1 for an uncaught exception, 3 for a syntax
// error (stderr then starts with "<line>:<col>: <message>").
package main

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/dop251/goja"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/parser"
)

//go:embed prelude.js
var prelude string

type exit struct{ code int }

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: js <file.js>")
		os.Exit(64)
	}
	src, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(66)
	}
	out := bufio.NewWriter(os.Stdout)
	code := run(string(src), out)
	_ = out.Flush()
	os.Exit(code)
}

func run(src string, out *bufio.Writer) (code int) {
	prog, err := goja.Compile(os.Args[1], src, false)
	if err != nil {
		var se *goja.CompilerSyntaxError
		if errors.As(err, &se) && se.File != nil {
			syntaxError(se.File, se.Offset, se.Message)
			return 3
		}
		var el parser.ErrorList
		if errors.As(err, &el) && len(el) > 0 && el[0].File != nil {
			syntaxError(el[0].File, el[0].Offset, el[0].Message)
			return 3
		}
		fmt.Fprintln(os.Stderr, err)
		return 3
	}
	vm := goja.New()
	vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
	stdin := func() string { b, _ := io.ReadAll(os.Stdin); return string(b) }
	host := map[string]any{
		"write": func(fd int, s string) {
			if fd == 2 {
				_ = out.Flush()
				fmt.Fprint(os.Stderr, s)
				return
			}
			_, _ = out.WriteString(s)
		},
		"stdin": func() string { return stdin() },
		"exit":  func(c int) { panic(exit{c}) },
	}
	if err := vm.Set("__host", host); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 70
	}
	if _, err := vm.RunString(prelude); err != nil {
		fmt.Fprintln(os.Stderr, "prelude:", err)
		return 70
	}
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(exit)
			if !ok {
				panic(r)
			}
			code = e.code
		}
	}()
	if _, err := vm.RunProgram(prog); err != nil {
		var ex *goja.Exception
		if errors.As(err, &ex) {
			if e, ok := ex.Value().Export().(exit); ok {
				return e.code
			}
			_ = out.Flush()
			fmt.Fprintln(os.Stderr, "Uncaught "+ex.String())
			return 1
		}
		var e *goja.InterruptedError
		if errors.As(err, &e) {
			fmt.Fprintln(os.Stderr, e)
			return 1
		}
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func syntaxError(f *file.File, offset int, msg string) {
	p := f.Position(offset)
	fmt.Fprintf(os.Stderr, "%d:%d: SyntaxError: %s\n", p.Line, p.Column, msg)
}
