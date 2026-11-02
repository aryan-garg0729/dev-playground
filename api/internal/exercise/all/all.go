// Package all registers every exercise type. Binaries import it for side effects.
package all

import (
	_ "github.com/aryangarg/codequest/api/internal/exercise/code" // code type
	_ "github.com/aryangarg/codequest/api/internal/exercise/mcq"  // mcq type
)
