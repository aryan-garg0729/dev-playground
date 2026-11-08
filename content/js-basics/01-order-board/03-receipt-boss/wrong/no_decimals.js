const input = require("fs").readFileSync(0, "utf8");
let total = 0;
for (const line of input.split("\n").filter((l) => l.trim() !== "")) {
  const [item, qty, price] = line.trim().split(" ");
  const cost = Number(qty) * Number(price);
  total += cost;
  console.log(`${item} x${qty} ${cost}`);
}
console.log(`TOTAL ${total}`);
