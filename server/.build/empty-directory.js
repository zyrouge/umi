const path = require("path");
const fs = require("fs");

const start = () => {
    const [dirPath] = process.argv.slice(2);
    if (!dirPath) {
        console.error("No directory path provided");
        process.exit(1);
    }
    const fullPath = path.resolve(dirPath);
    if (fs.existsSync(fullPath)) {
        fs.rmSync(fullPath, { recursive: true });
    }
};

start();
