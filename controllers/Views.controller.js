const { logRequest, logError } = require("./Controller.js");
const { Views } = require("../models/index.js");

exports.test = (req, res) => {
  console.log("Testing... It's Working!!!");
  res.send({ message: "OK" });
};

exports.findAll = (req, res) => {
  let condition = Object.keys(req.query).length ? req.query : null;
  logRequest("GET", "TABLE = 'Views' CONDITION =", condition);

  Views.findAll({ where: condition })
    .then((data) => res.send(data))
    .catch((err) => logError(err) || res.status(500).send({ message: err.message }));
};

exports.findOne = (req, res) => {
  // logRequest("GET", "TABLE = 'Views' ID =", req.params.id);

  Views.findAll({ where: { id: req.params.id } })
    .then((data) => res.send(data))
    .catch((err) => logError(err) || res.status(500).send({ message: err.message }));
};

exports.create = (req, res) => {
  logRequest("POST", "TABLE = 'Views' DATA =", req.body);

  let { Curr_Date, Count } = req.body;

  if (!Curr_Date || !Count) {
    logError("Invalid request message parameters");
    res.status(400).send({ message: "Invalid request message parameters" });
    return;
  }

  // Create new structure
  let newViews = { Curr_Date: Curr_Date, Count: Count };

  Views.create(newViews)
    .then((data) => res.status(201).send(data))
    .catch((err) => logError(err) || res.status(500).send({ message: err.message }));
};

exports.update = (req, res) => {
  logRequest("PUT", "TABLE = 'Views' CONDITION =", { ...req.query, ...req.params }, "DATA = ", req.body);

  // Check if parameters and updated data
  if (!Object.keys(req.body).length || (!Object.keys(req.query).length && !Object.keys(req.params).length)) {
    logError("Invalid request message parameters");
    res.status(400).send({ message: "Invalid request message parameters" });
    return;
  }

  Views.update(req.body, { where: { ...req.query, ...req.params } })
    .then((data) => (data ? res.sendStatus(204) : res.sendStatus(304)))
    .catch((err) => logError(err) || res.status(500).send({ message: err }));
};

exports.destroy = (req, res) => {
  logRequest("DELETE", "TABLE = 'Views' CONDITION =", { ...req.query, ...req.params });

  if (!Object.keys(req.query).length && !Object.keys(req.params).length) {
    logError("Invalid request message parameters");
    res.status(400).send({ message: "Invalid request message parameters" });
    return;
  }

  Views.destroy({ where: { ...req.query, ...req.params } })
    .then((data) => (data ? res.sendStatus(204) : res.sendStatus(304)))
    .catch((err) => logError(err) || res.status(500).send({ message: err }));
};
