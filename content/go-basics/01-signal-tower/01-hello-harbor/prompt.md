Every Go program starts in `package main`, in a function called `main`. When `main` returns, the program is done.

To print a line, import the `fmt` package (short for "format") and call `fmt.Println`. It prints its arguments and a newline.

```go
package main

import "fmt"

func main() {
	fmt.Println("Lamp check")
}
// Lamp check
```

## Assignment

Welcome to the night shift at Saltmere harbor. Mira wants the signal board to announce that the tower is up.

Make the program print exactly:

```
Saltmere tower online
```
