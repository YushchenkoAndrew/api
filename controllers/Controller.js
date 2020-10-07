exports.logRequest = (...info) => {
  console.log(`\x1b[34m[${info[0]} REQUEST]\x1b[0m\x1b[35m[${getTime()}]\x1b[0m`, ...info.slice(1));
  function getTime() {
    let arr = new Date().toString().split(" ");
    return arr[4] + " " + [arr[2], arr[1], arr[3]].join("-");
  }
};

exports.logError = (err) => console.log(`\x1b[31m[Error]\x1b[0m ${err}`);
