// Command go is the Go interpreter (yaegi) compiled to WASI. The grading host
// runs it as `go /work/<entry>`; the program's stdin, stdout and stderr are the
// learner program's.
//
// Exit codes: the program's own, 2 for a runtime panic, 3 for a compile error
// (stderr then starts with "<line>:<col>: <message>").
package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/traefik/yaegi/interp"

	"github.com/aryangarg/codequest/runners/wasm/symbols"
)

var (
	posRe    = regexp.MustCompile(`:(\d+):(\d+): (.*)`)
	importRe = regexp.MustCompile(`^import "([^"]+)" error: .*`)
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go <file.go>")
		os.Exit(64)
	}
	i := interp.New(interp.Options{GoPath: "/nonexistent", Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr})
	if err := i.Use(symbols.Symbols); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(70)
	}
	prog, err := i.CompilePath(os.Args[1])
	if err != nil {
		msg := err.Error()
		if m := posRe.FindStringSubmatch(msg); m != nil {
			msg = m[1] + ":" + m[2] + ": " + importRe.ReplaceAllString(m[3], `package "$1" is not available in lessons`)
		}
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(3)
	}
	if _, err := i.Execute(prog); err != nil {
		var p interp.Panic
		if errors.As(err, &p) {
			fmt.Fprintf(os.Stderr, "panic: %v\n", p.Value)
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(2)
	}
}
