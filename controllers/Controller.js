const { logRequest, logDebug } = require("../lib/log");
const { errorHandler, resByType, getDataByType } = require("../lib/resHandler");
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
exports.findAll = (req, res, next) => {
  let table = getTable(req.params.table);
  let query = Object.keys(req.query).length ? req.query : null;
  logRequest("GET", `TABLE = '${table}' QUERY =`, query);

  if (table)
    table.find(
      query,
      // Callback
      (data) => (data.length ? resByType(req.header("Content-Type"), data, res) : res.sendStatus(204)),
      // Error
      (err) => errorHandler(500, err.message, req, res)
    );
  else next();
};

// REDIRECT GET Request to the <Table>
exports.findOne = (req, res, next) => {
  let table = getTable(req.params.table);
  logRequest("GET", `TABLE = '${table}' ID =`, req.params.id);

  // Bad Request handler
  if (req.params.id === undefined || isNaN(Number(req.params.id))) {
    errorHandler(400, "Invalid request message parameters", req, res);
    return;
  }

  if (table)
    table.find(
      // Query
      { id: req.params.id },
      // Callback
      (data) => (data.length ? resByType(req.header("Content-Type"), data, res) : res.sendStatus(204)),
      // Error
      (err) => errorHandler(500, err.message, req, res)
    );
  else next();
};

// REDIRECT POST Request to the <Table>
exports.create = (req, res, next) => {
  let table = getTable(req.params.table);
  // delete req.params.table;
  logRequest("POST", `TABLE = '${table}' DATA =`, req.body);

  if (table) {
    let data = new table.Data(getDataByType(req.header("Content-Type"), req.body));

    // Bad Request handler
    if (data.check()) {
      errorHandler(400, "Invalid request message parameters", req, res);
      return;
    }

    table.create(
      data,
      // Callback
      (data) => resByType(req.header("Content-Type"), data, res, 201),
      // Error
      (err) => errorHandler(500, err.message, req, res)
    );
  } else next();
};

// REDIRECT PUT Request to the <Table>
exports.update = (req, res, next) => {
  let table = getTable(req.params.table);
  delete req.params.table;

  logRequest("PUT", `TABLE = '${table}' CONDITION =`, { ...req.query, ...req.params }, "DATA = ", req.body);

  // Check if parameters and updated data
  if (!Object.keys(req.body).length || !(!Object.keys(req.query).length ^ !Object.keys(req.params).length)) {
    errorHandler(400, "Invalid request message parameters", req, res);
    return;
  }

  if (table)
    table.update(
      // Data
      getDataByType(req.header("Content-Type"), req.body),
      // Query
      { ...req.query, ...req.params },
      // Callback
      ([data]) => (data ? res.sendStatus(204) : res.sendStatus(304)),
      // Error
      (err) => errorHandler(500, err.message, req, res)
    );
  else next();
};

// REDIRECT DELETE Request to the <Table>
exports.destroy = (req, res, next) => {
  let table = getTable(req.params.table);
  delete req.params.table;
  logRequest("DELETE", `TABLE = '${table}' CONDITION =`, { ...req.query, ...req.params });

  // Bad Request handler
  if (!(!Object.keys(req.query).length ^ !Object.keys(req.params).length)) {
    errorHandler(400, "Invalid request message parameters", req, res);
    return;
  }

  if (table)
    table.destroy(
      // Query
      { ...req.query, ...req.params },
      // Callback
      (data) => (data ? res.sendStatus(204) : res.sendStatus(304)),
      // Error
      (err) => errorHandler(500, err.message, req, res)
    );
  else next();
};
