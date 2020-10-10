// Basic Output Style
exports.logRequest = (req, ...info) => {
  console.log(`\x1b[34m[${req} REQUEST]\x1b[0m\x1b[35m[${getTime()}]\x1b[0m`, ...info);
  function getTime() {
    let arr = new Date().toString().split(" ");
    return arr[4] + " " + [arr[2], arr[1], arr[3]].join("-");
  }
};

// Basic Output Style
exports.logError = (...err) => console.log("\x1b[31m[Error]\x1b[0m", ...err);

// Basic Output Style
exports.logInfo = (...message) => console.log("\x1b[32m[INFO]\x1b[0m", ...message);
