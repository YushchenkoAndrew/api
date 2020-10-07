const router = require("express").Router();
const visitor = require("../controllers/Visitor.controller");

//
// Table "Visitor"
//

// Return all or requested one data from Table "Visitor"
router.get("/Visitor", visitor.findAll);
// Return a single data row from "Visitor" with id
router.get("/Visitor/:id", visitor.findOne);
// Create new data into "Visitor"
router.post("/Visitor", visitor.create);
// Update data in "Visitor" found by query
router.put("/Visitor", visitor.update);
// Update data in "Visitor" by id
router.put("/Visitor/:id", visitor.update);
// Delete data from "Visitor" that found by query
router.delete("/Visitor", visitor.delete);
// Delete data from "Visitor" by id
router.delete("/Visitor/:id", visitor.delete);

module.exports = router;
