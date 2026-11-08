package main

import "fmt"

type ship struct {
	name    string
	tonnage int
}

func main() {
	var fleet []ship
	fmt.Println(len(fleet))
}
