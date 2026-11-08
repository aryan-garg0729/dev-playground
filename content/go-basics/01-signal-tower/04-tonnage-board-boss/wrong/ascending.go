package main

import (
	"fmt"
	"sort"
)

type ship struct {
	name    string
	tonnage int
}

func main() {
	var fleet []ship
	for {
		var s ship
		if _, err := fmt.Scan(&s.name, &s.tonnage); err != nil {
			break
		}
		fleet = append(fleet, s)
	}
	sort.SliceStable(fleet, func(i, j int) bool { return fleet[i].tonnage < fleet[j].tonnage })
	total := 0
	for _, s := range fleet {
		fmt.Println(s.name, s.tonnage)
		total += s.tonnage
	}
	fmt.Println("total", total)
}
