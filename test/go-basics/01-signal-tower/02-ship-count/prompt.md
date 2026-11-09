A program that always prints the same thing gets boring fast. Real programs read input. The tower's teleprinter types into your program's standard input (stdin).

`fmt.Scan` reads the next value from stdin into a variable. You pass it the variable's address with `&`, so it can fill it in.

```go
var lamps int
fmt.Scan(&lamps)
fmt.Println(lamps, "lamps lit")
// input: 4
// 4 lamps lit
```

`fmt.Println` puts a space between its arguments. If you need more control, `fmt.Printf("%d lamps lit\n", lamps)` formats the line yourself. `%d` is a whole number.

## Assignment

Tomas is tired of counting hulls by hand. Read one whole number, the ships in the harbor, and print how many are docked:

- `3` prints `3 ships docked`
- `1` prints `1 ship docked` (one ship has no "s")

Every other count, zero included, uses `ships`.
