require('shelljs/global');
config.fatal = true;

const { desc, task, namespace } = require("jake");
const sh = require("shelljs");
const { hasNewerFile } = require("./jakelib/lib/incremental.cjs");

desc("Build");
task("build", async () => {
    if (! hasNewerFile(["./src"], "build")) {
        console.log("No update")
        return;
    }
    sh.exec("craco build");
});

desc("Generate protocol-buffer connector");
task("bufgen", async () => {
    const sh = require("shelljs");
    sh.mkdir("-p", "src/pbgen");
    sh.exec("buf generate ../proto/");
});

desc("Start development server");
task("start", async () => {
    sh.exec("craco start");
});

desc("Generate all");
task("gen", ["bufgen"]);
