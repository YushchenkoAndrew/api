const express = require("express");
const bodyParser = require("body-parser");
const app = express();
const router = require("./routes/routes");
const db = require("./models/index");
const swaggerUi = require("swagger-ui-express");
const swaggerDocument = require("./doc/swagger.json");

const HOST = "0.0.0.0";
const PORT = 31337;

// parse requests of content-type - application/json
app.use(bodyParser.json());

// parse requests of content-type - application/x-www-form-urlencoded
app.use(bodyParser.urlencoded({ extended: true }));

// Swagger Documentation
app.use("/api/doc", swaggerUi.serve, swaggerUi.setup(swaggerDocument));

// Determining Routers settings
app.use("/api", router);

process.on("SIGINT", async function () {
  console.log("\x1b[32m[INFO]\x1b[0m Server Terminated");
  await db.sequelize.close();
  process.exit(0);
});

app.listen(PORT, HOST, (err) => {
  if (err) console.log(`\x1b[31m[ERROR]\x1b[0m Error appear ${err}`);
  console.log("\x1b[32m[INFO]\x1b[0m Server Started ...");
  console.log(`\x1b[32m[INFO]\x1b[0m Listening on Port ${PORT}`);
});
