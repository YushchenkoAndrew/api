const router = require("express").Router();
const { findAll, findOne, create, update, destroy } = require("../controllers/Controller");
const { authorizationToken: auth } = require("../middleware/auth");

// DataBase
// destroy == DELETE Request
// Return all or requested one data from <Table>
router.get("/:table", findAll);
// Return a single data row from <Table> with id
router.get("/:table/:id", findOne);
// Create new data into "<Table>
router.post("/:table", auth, create);
// Update data in <Table> found by query
router.put("/:table", auth, update);
// Update data in <Table> by id
router.put("/:table/:id", auth, update);
// Delete data from <Table> that found by query
router.delete("/:table", auth, destroy);
// Delete data from <Table> by id
router.delete("/:table/:id", auth, destroy);

module.exports = router;
