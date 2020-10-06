const express = require("express");
const bodyParser = require("body-parser");
const app = express();
const router = require("./config/routes");
const db = require("./models/index");

const HOST = "0.0.0.0";
const PORT = 31337;

// parse requests of content-type - application/json
app.use(bodyParser.json());

// parse requests of content-type - application/x-www-form-urlencoded
app.use(bodyParser.urlencoded({ extended: true }));

// Determining Routers settings
app.use("/api", router);

process.on("SIGINT", function () {
  console.log("\x1b[32m[INFO]\x1b[0m Server Terminated");
  process.exit();
});

app.listen(PORT, HOST, (err) => {
  if (err) console.log(`\x1b[31m[ERROR]\x1b[0m Error appear ${err}`);
  console.log("\x1b[32m[INFO]\x1b[0m Server Started ...");
  console.log(`\x1b[32m[INFO]\x1b[0m Listening on Port ${PORT}`);
});
