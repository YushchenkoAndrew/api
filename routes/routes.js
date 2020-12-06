const router = require("express").Router();
const dbTables = require("./dbTables.routes.js");
const { logRequest } = require("../lib/log.js");
const { authorizationToken: auth, generateToken } = require("../middleware/auth");

// Check if API is Alive
router.get("/ping", (req, res) => logRequest("GET", "PING") || res.send({ message: "OK" }));

// Authorization
router.get("/token", auth, (req, res) => res.send({ message: "OK" }));
router.post("/login", generateToken);

// DataBase Tables
router.use("/", dbTables);

module.exports = router;
