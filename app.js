const { logInfo, logError } = require("./lib/log");
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

// Catch 404 and forward to error handler
app.use((req, res, next) => logError("Not Found", req.url) || res.status(404).send({ message: "Not Found" }));

process.on("SIGINT", async function () {
  logInfo("Server Terminated");
  await db.sequelize.close();
  process.exit(0);
});

app.listen(PORT, HOST, (err) => {
  if (err) logError(`Error appear ${err}`);
  logInfo("Server Started ...");
  logInfo(`Listening on Port ${PORT}`);
});
