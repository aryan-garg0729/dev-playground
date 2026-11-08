package main

import "fmt"

func main() {
	var ships int
	fmt.Scan(&ships)
	if ships == 1 {
		fmt.Println("1 ship docked")
		return
	}
	fmt.Printf("%d ships docked\n", ships)
}
