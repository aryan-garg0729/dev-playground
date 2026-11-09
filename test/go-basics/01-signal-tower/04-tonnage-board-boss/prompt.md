Mira's ticket: *"The harbor pilots need the heaviest ships first. Put the night's arrivals on the board, biggest first, with the total at the bottom."*

You have everything you need from this chapter, plus two new tools.

**Slices of structs.** A struct groups fields; a slice grows with `append`:

```go
type ship struct {
	name    string
	tonnage int
}

var fleet []ship
fleet = append(fleet, ship{"kestrel", 1200})
```

**Reading until the input ends.** `fmt.Scan` returns how many values it read and an error. At the end of the input the error is not `nil`, so you can loop until then:

```go
for {
	var s ship
	if _, err := fmt.Scan(&s.name, &s.tonnage); err != nil {
		break
	}
	fleet = append(fleet, s)
}
```

**Sorting.** `sort.SliceStable(fleet, less)` sorts with your `less(i, j)` function and keeps equal elements in their original order.

## Assignment

Read every line of stdin. Each is a ship name and its tonnage: `kestrel 1200`.

1. Print each ship as `<name> <tonnage>`, heaviest first. Ships with the same tonnage stay in the order they arrived.
2. Then print `total <sum of tonnages>`.

With no ships, print just `total 0`.

```
input:          output:
kestrel 1200    osprey 3000
marlin 800      kestrel 1200
osprey 3000     marlin 800
                total 5000
```
