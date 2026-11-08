const rolls = Number(require("fs").readFileSync(0, "utf8").trim());
console.log(`${Math.floor(rolls / 12)} boxes, ${rolls % 12} loose`);
