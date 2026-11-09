Pia's ticket: *"The till broke. Customers' orders come in as lines; print a receipt with each line's price and a total."*

New tools for the job:

```js
"rye 2 4.50".split(" ");      // [ 'rye', '2', '4.50' ]
"a\n\nb".split("\n");        // [ 'a', '', 'b' ]
[1, 0, 3].filter((x) => x);   // [ 1, 3 ]
(9).toFixed(2);               // '9.00'
```

`toFixed(2)` turns a number into a string with exactly two decimals.

Money and floating point don't always agree (`0.1 * 3` is `0.30000000000000004`). Here, `toFixed(2)` on each line's price and on the total rounds that away. Real tills count in cents. We won't go further than that here.

## Assignment

Read all of stdin. Each non-blank line is `<item> <quantity> <unit price>`.

1. For each order line, in order, print `<item> x<quantity> <line price>`, where the line price is quantity × unit price with two decimals.
2. Then print `TOTAL <sum of line prices>`, with two decimals.

Skip blank lines. With no orders, print just `TOTAL 0.00`.

```
input:              output:
rye 2 4.50          rye x2 9.00
croissant 3 2.25    croissant x3 6.75
                    TOTAL 15.75
```
