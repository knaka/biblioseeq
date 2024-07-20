const sh = require("shelljs");

const getLatestModifiedTime = (dir) => {
    const fs = require('fs');
    const path = require('path');

    let latestTime = 0;
    const files = fs.readdirSync(dir);
    files.forEach((file) => {
        const filePath = path.join(dir, file);
        const stats = fs.statSync(filePath);
        if (stats.isDirectory()) {
            const innerTime = getLatestModifiedTime(filePath);
            if (innerTime > latestTime) {
                latestTime = innerTime;
            }
        } else {
            if (stats.mtimeMs > latestTime) {
                latestTime = stats.mtimeMs
            }
        }
    });

    return latestTime;
};

// Check if any of the source directories have a newer file than the target directory
exports.hasNewerFile = (srcDirs, targetDir) => {
    sh.mkdir("-p", targetDir);
    // If the target directory is empty, we need to build
    if (sh.ls(targetDir).length === 0) {
        return true;
    }
    const latestBuildTime = getLatestModifiedTime(targetDir);
    return srcDirs.some((dir) => getLatestModifiedTime(dir) > latestBuildTime);
}
