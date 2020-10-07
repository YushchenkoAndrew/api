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
//
router.put("/Visitor", visitor.update);
router.put("/Visitor/:id", visitor.update);

module.exports = router;
