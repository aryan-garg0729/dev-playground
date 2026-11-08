package symbols

import (
	"errors"
	"os"
)

// yaegi extract wraps these two os functions (see the generated os.go, which
// linters skip). os.Exit really exits, so the host reads the code; there are no
// other processes to find under WASI.

//nolint:unused // used by the generated os.go
func osExit(code int) { os.Exit(code) }

//nolint:unused // used by the generated os.go
func osFindProcess(int) (*os.Process, error) { return nil, errors.New("os: no processes") }
