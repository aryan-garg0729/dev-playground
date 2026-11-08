const rolls = Number(require("fs").readFileSync(0, "utf8").trim());
console.log(`${Math.floor(rolls / 10)} boxes, ${rolls % 10} loose`);
