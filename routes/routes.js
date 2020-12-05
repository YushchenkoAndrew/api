const router = require("express").Router();
const { findAll, findOne, create, update, destroy } = require("../controllers/Controller");
const { authorizationToken: auth, generateToken } = require("../middleware/auth");
// destroy == DELETE Request

// Return all or requested one data from <Table>
router.get("/:table", findAll);
// Return a single data row from <Table> with id
router.get("/:table/:id", findOne);
// Create new data into "<Table>
router.post("/:table", create);
// Update data in <Table> found by query
router.put("/:table", update);
// Update data in <Table> by id
router.put("/:table/:id", update);
// Delete data from <Table> that found by query
router.delete("/:table", destroy);
// Delete data from <Table> by id
router.delete("/:table/:id", destroy);

// Check if API is Alive
router.get("/ping", (req, res) => res.json({ message: "OK" }));

// Authorization
router.get("/ping2", auth, (req, res) => res.json({ message: "OK" }));
router.post("/login", generateToken);

module.exports = router;
