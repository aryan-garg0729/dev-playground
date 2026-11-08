The board's input arrives on standard input (stdin). Read all of it as one string:

```js
const input = require("fs").readFileSync(0, "utf8");
```

That string still has the newline at the end, and it's text, not a number. `Number(input.trim())` fixes both.

Two pieces of arithmetic come up every day in a bakery:

```js
Math.floor(26 / 12); // 2   (whole boxes; floor drops the fraction)
26 % 12;             // 2   (the remainder: what's left over)
```

## Assignment

Bram bakes rolls and packs them in boxes of 12. Read the number of rolls from stdin and print how many full boxes he can pack and how many rolls are left loose:

```
26  ->  2 boxes, 2 loose
```

Keep the wording exactly, even for `1 boxes`. Bram says grammar is the counter staff's job.
