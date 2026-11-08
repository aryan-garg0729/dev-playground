const input = require("fs").readFileSync(0, "utf8");
let total = 0;
for (const line of input.trim().split("\n")) {
  const [item, qty, price] = line.split(" ");
  const cost = Number(qty) * Number(price);
  total += cost;
  console.log(`${item} x${qty} ${cost.toFixed(2)}`);
}
console.log(`TOTAL ${total.toFixed(2)}`);
