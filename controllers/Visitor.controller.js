const { logRequest, logError } = require("./Controller");
const { Visitors } = require("../models/index.js");

exports.test = (req, res) => {
  console.log("Testing... It's Working!!!");
  res.send({ message: "OK" });
};

exports.findAll = (req, res) => {
  let condition = Object.keys(req.query).length ? req.query : null;
  logRequest("GET", "TABLE = 'Visitor' CONDITION =", condition);

  Visitors.findAll({ where: condition })
    .then((data) => res.send(data))
    .catch((err) => logError(err) || res.status(500).send({ message: err.message }));
};

exports.findOne = (req, res) => {
  logRequest("GET", "TABLE = 'Visitor' ID =", req.params.id);

  Visitors.findAll({ where: { id: req.params.id } })
    .then((data) => res.send(data))
    .catch((err) => logError(err) || res.status(500).send({ message: err.message }));
};

exports.create = (req, res) => {
  logRequest("POST", "TABLE = 'Visitor' DATA =", req.body);

  let { Country, ip, Visit_Date, Count } = req.body;

  if (!Country || !ip || !Visit_Date || !Count) {
    logError("Invalid request message parameters");
    res.status(400).send({ message: "Invalid request message parameters" });
    return;
  }

  // Create new structure
  let newVisitor = { Country: Country, ip: ip, Visit_Date: Visit_Date, Count: Count };

  Visitors.create(newVisitor)
    .then((data) => res.status(200).send(data))
    .catch((err) => logError(err) || res.status(500).send({ message: err.message }));
};

exports.update = (req, res) => {
  logRequest("PUT", "TABLE = 'Visitor' CONDITION =", { ...req.query, ...req.params }, "DATA = ", req.body);

  // Check if parameters and updated data
  if (!Object.keys(req.body).length || (!Object.keys(req.query).length && !Object.keys(req.params).length)) {
    logError("Invalid request message parameters");
    res.status(400).send({ message: "Invalid request message parameters" });
    return;
  }

  Visitors.update(req.body, { where: { ...req.query, ...req.params } })
    .then((data) => (data ? res.send({ message: "Information was updated successfully" }) : res.send({ message: "Information updating was Failed" })))
    .catch((err) => logError(err) || res.status(500).send({ message: err }));
};

exports.destroy = (req, res) => {
  logRequest("DELETE", "TABLE = 'Visitor' CONDITION =", { ...req.query, ...req.params });

  if (!Object.keys(req.query).length && !Object.keys(req.params).length) {
    logError("Invalid request message parameters");
    res.status(400).send({ message: "Invalid request message parameters" });
    return;
  }

  Visitors.destroy({ where: { ...req.query, ...req.params } })
    .then((data) => (data ? res.send({ message: "Information was deleted successfully" }) : res.send({ message: "Information deleting was Failed" })))
    .catch((err) => logError(err) || res.status(500).send({ message: err }));
};
