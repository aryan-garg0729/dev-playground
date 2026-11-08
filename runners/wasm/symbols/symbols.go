// Package symbols holds the standard library a Go lesson may import, generated
// by `make symbols` (yaegi extract). Packages not listed fail to compile, which
// keeps the module small; WASI already blocks the network and every file
// outside /work.
package symbols

import "reflect"

// Symbols is filled by the generated files' init functions.
var Symbols = map[string]map[string]reflect.Value{}
