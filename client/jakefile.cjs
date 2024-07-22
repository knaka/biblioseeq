require('shelljs/global');
config.fatal = true;

const { desc, task, namespace } = require("jake");
const sh = require("shelljs");
const { hasNewerFile } = require("./jakelib/lib/incremental.cjs");

desc("Build for production");
const buildTask = task("build", async () => {
    sh.exec("craco build");
});

namespace("build", () => {
    desc("Build for development")
    task("dev", async () => {
        if (! hasNewerFile(["./src"], "build")) {
            return;
        }
        sh.exec("craco build");
    });

    desc(buildTask.description);
    task("prd", buildTask.prereqs, buildTask.action);
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
