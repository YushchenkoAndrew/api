const { logRequest, logError } = require("../lib/log");
const visitor = require("../services/Visitor.service");
const views = require("../services/Views.service");
const github = require("../services/Github.service");
const usedLanguages = require("../services/UsedLanguages.service");
//
// REDIRECTION
//

function getTable(name) {
  switch (name) {
    case "Visitor":
      return visitor;

    case "Views":
      return views;

    case "Github":
      return github;

    case "UsedLanguages":
      return usedLanguages;
  }
}

// REDIRECT GET Request to the <Table>
exports.findAll = (req, res) => {
  let table = getTable(req.params.table);
  let query = Object.keys(req.query).length ? req.query : null;
  logRequest("GET", `TABLE = '${table}' QUERY =`, query);

  if (table)
    table.find(
      query,
      // Callback
      (data) => res.send(data),
      // Error
      (err) => logError(err) || res.status(500).send({ message: err.message })
    );
  else logError("Table not Found") || res.status(404).send({ message: "Table not Found" });
};

// REDIRECT GET Request to the <Table>
exports.findOne = (req, res) => {
  let table = getTable(req.params.table);
  logRequest("GET", `TABLE = '${table}' ID =`, req.params.id);

  // Bad Request handler
  if (req.params.id === undefined || isNaN(Number(req.params.id))) {
    logError("Invalid request message parameters");
    res.status(400).send({ message: "Invalid request message parameters" });
    return;
  }

  if (table)
    table.find(
      // Query
      { id: req.params.id },
      // Callback
      (data) => res.send(data),
      // Error
      (err) => logError(err) || res.status(500).send({ message: err.message })
    );
  else logError("Table not Found") || res.status(404).send({ message: "Table not Found" });
};

// REDIRECT POST Request to the <Table>
exports.create = (req, res) => {
  let table = getTable(req.params.table);
  // delete req.params.table;
  logRequest("POST", `TABLE = '${table}' DATA =`, req.body);

  if (table) {
    let data = new table.Data(req.body);

    // Bad Request handler
    if (data.check()) {
      logError("Invalid request message parameters");
      res.status(400).send({ message: "Invalid request message parameters" });
      return;
    }

    table.create(
      data,
      // Callback
      (data) => res.status(201).send(data),
      // Error
      (err) => logError(err) || res.status(500).send({ message: err.message })
    );
  } else logError("Table not Found") || res.status(404).send({ message: "Table not Found" });
};

// REDIRECT PUT Request to the <Table>
exports.update = (req, res) => {
  let table = getTable(req.params.table);
  delete req.params.table;

  logRequest("PUT", `TABLE = '${table}' CONDITION =`, { ...req.query, ...req.params }, "DATA = ", req.body);

  // Check if parameters and updated data
  if (!Object.keys(req.body).length || !(!Object.keys(req.query).length ^ !Object.keys(req.params).length)) {
    logError("Invalid request message parameters");
    res.status(400).send({ message: "Invalid request message parameters" });
    return;
  }

  if (table)
    table.update(
      // Data
      req.body,
      // Query
      { ...req.query, ...req.params },
      // Callback
      ([data]) => (data ? res.sendStatus(204) : res.sendStatus(304)),
      // Error
      (err) => logError(err) || res.status(500).send({ message: err.message })
    );
  else logError("Table not Found") || res.status(404).send({ message: "Table not Found" });
};

// REDIRECT DELETE Request to the <Table>
exports.destroy = (req, res) => {
  let table = getTable(req.params.table);
  delete req.params.table;
  logRequest("DELETE", `TABLE = '${table}' CONDITION =`, { ...req.query, ...req.params });

  // Bad Request handler
  if (!(!Object.keys(req.query).length ^ !Object.keys(req.params).length)) {
    logError("Invalid request message parameters");
    res.status(400).send({ message: "Invalid request message parameters" });
    return;
  }

  if (table)
    table.destroy(
      // Query
      { ...req.query, ...req.params },
      // Callback
      (data) => (data ? res.sendStatus(204) : res.sendStatus(304)),
      // Error
      (err) => logError(err) || res.status(500).send({ message: err })
    );
  else logError("Table not Found") || res.status(404).send({ message: "Table not Found" });
};
