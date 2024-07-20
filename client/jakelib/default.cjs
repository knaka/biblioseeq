const {desc, task, run} = require("jake");

desc("Display the tasks with descriptions");
task("default", function () {
    const originalConsoleLog = console.log;
    console.log = function(...args) {
        originalConsoleLog((args.join(" ")).replace(/^jake[[ \t]*/, ""));
    };
    const { run } = require("jake");
    run.apply(global.jake, ["--tasks"]);
});
