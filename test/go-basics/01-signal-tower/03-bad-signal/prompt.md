When a program finishes, it hands the system a number: its exit code. `0` means "all good". Anything else means something went wrong, and scripts that run your program check it.

Returning from `main` exits with `0`. To exit with another code, call `os.Exit`:

```go
import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "lamp broken")
	os.Exit(1)
}
// (on stderr) lamp broken
// exit code: 1
```

`fmt.Fprintln(os.Stderr, ...)` prints to standard error, the channel for problems, so they don't get mixed into the real output.

`switch` picks one branch, and a case may list several values:

```go
switch tide {
case "high", "low":
	fmt.Println("normal")
default:
	fmt.Println("odd tide")
}
```

## Assignment

A gull sat on the teleprinter last night, and the board lit up with nonsense. Read one word, a signal code, with `fmt.Scan`.

- For `red`, `green` or `amber`, print `OK <code>`, for example `OK green`, and exit normally.
- For anything else, including no input at all, print `unknown signal: <code>` to **stderr** and exit with code `2`.
