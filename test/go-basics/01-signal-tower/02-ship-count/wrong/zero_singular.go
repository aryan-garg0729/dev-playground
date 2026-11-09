package main

import "fmt"

func main() {
	var ships int
	fmt.Scan(&ships)
	if ships <= 1 {
		fmt.Printf("%d ship docked\n", ships)
		return
	}
	fmt.Printf("%d ships docked\n", ships)
}
