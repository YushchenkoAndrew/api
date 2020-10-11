const { logInfo, logDebug } = require("./lib/log");
const { errorHandler } = require("./lib/errorHandler");
const express = require("express");
const bodyParser = require("body-parser");
require("body-parser-xml")(bodyParser);
const app = express();
const router = require("./routes/routes");
const db = require("./models/index");
const swaggerUi = require("swagger-ui-express");
const swaggerDocument = require("./doc/swagger.json");

const HOST = "0.0.0.0";
const PORT = 31337;

// parse requests of content-type - application/json
app.use(bodyParser.json());

// parse requests of content-type - application/xml
app.use(bodyParser.xml());

// parse requests of content-type - application/x-www-form-urlencoded
app.use(bodyParser.urlencoded({ extended: true }));

// FIXME: Fix PUT REQUEST in swagger.json + Add XML REQUEST HANDLER
// Swagger Documentation
app.use("/api/doc", swaggerUi.serve, swaggerUi.setup(swaggerDocument));

// Determining Routers settings
app.use("/api", router);

// Catch 404 and forward to error handler
app.use((req, res, next) => errorHandler(404, `Not Found '${req.url}'`, req, res));

// Error handler
app.use((err, req, res, next) => errorHandler(500, `Not Found '${req.url}'`, req, res));

process.on("SIGINT", async function () {
  logInfo("Server Terminated");
  await db.sequelize.close();
  process.exit(0);
});

app.listen(PORT, HOST, (err) => {
  if (err) errorHandler(500, `Error appear ${err.message}`);
  logInfo("Server Started ...");
  logInfo(`Listening on Port ${PORT}`);
});
