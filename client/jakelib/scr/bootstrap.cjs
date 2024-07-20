const fs = require("fs");
const path = require("path");
const jakelibDir = path.join(__dirname, "..");

require("../../jakefile.cjs");

const files = fs.readdirSync(jakelibDir).filter(file => file.endsWith(".cjs"));
files.forEach(file => {
  const filePath = path.join(jakelibDir, file);
  require(filePath);
});
