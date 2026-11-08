package main

import (
	"fmt"
	"os"
)

func main() {
	var code string
	fmt.Scan(&code)
	switch code {
	case "red", "green", "amber":
		fmt.Println("OK " + code)
	default:
		fmt.Fprintln(os.Stderr, "unknown signal: "+code)
		os.Exit(2)
	}
}
