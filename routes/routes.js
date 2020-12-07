const router = require("express").Router();
const dbTables = require("./dbTables.routes.js");
const { logRequest } = require("../lib/log.js");
const { menu, getTables } = require("../controllers/menu.js");
const { authorizationToken: auth, generateToken } = require("../middleware/auth");

// API Navigation
router.get("/", menu);
router.get("/tables", getTables);

// Check if API is Alive
router.get("/ping", (req, res) => logRequest("GET", "PING") || res.send({ message: "OK" }));

// Authorization
router.get("/token", auth, (req, res) => res.send({ message: "OK" }));
router.post("/login", generateToken);

// DataBase Tables
router.use("/", dbTables);

module.exports = router;
