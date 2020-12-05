const { logInfo, logDebug } = require("./lib/log");
const { errorHandler } = require("./lib/resHandler");
const express = require("express");
const bodyParser = require("body-parser");
require("body-parser-xml")(bodyParser);
const app = express();
const router = require("./routes/routes");
const db = require("./models/index");
const swaggerUi = require("swagger-ui-express");
const swaggerDocument = require("./doc/swagger.json");
require("dotenv").config();

const HOST = process.env.API_HOST || "127.0.0.1";
const PORT = process.env.API_PORT || 31337;

// parse requests of content-type - application/json
app.use(bodyParser.json());

// parse requests of content-type - application/xml
app.use(bodyParser.xml());
// TODO: think how to send xml file similar to json (POST REQUEST) + how to send it back

// parse requests of content-type - application/x-www-form-urlencoded
app.use(bodyParser.urlencoded({ extended: true }));

// FIXME: Fix PUT REQUEST in swagger.json + Add XML REQUEST HANDLER
// FIXME: + add httpStatusCode in GET REQUEST (STATUS = 204 "NO CONTENT")
// Swagger Documentation
app.use("/api/doc", swaggerUi.serve, swaggerUi.setup(swaggerDocument));

// Determining Routers settings
app.use("/api", router);

// Catch 404 and forward to error handler
app.use((req, res, next) => errorHandler(404, `Not Found '${req.url}'`, req, res));

// Error handler
app.use((err, req, res, next) => errorHandler(err.status, err.message.split("\n")[0], req, res));

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
