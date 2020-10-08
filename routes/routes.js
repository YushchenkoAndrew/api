const router = require("express").Router();
const { findAll, findOne, create, update, destroy } = require("../controllers/Controller");
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

module.exports = router;
