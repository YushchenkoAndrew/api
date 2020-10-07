// Basic Output Style
exports.logRequest = (...info) => {
  console.log(`\x1b[34m[${info[0]} REQUEST]\x1b[0m\x1b[35m[${getTime()}]\x1b[0m`, ...info.slice(1));
  function getTime() {
    let arr = new Date().toString().split(" ");
    return arr[4] + " " + [arr[2], arr[1], arr[3]].join("-");
  }
};

// Basic Output Style
exports.logError = (err) => console.log(`\x1b[31m[Error]\x1b[0m ${err}`);

//
// REDIRECTION
//
const visitor = require("./Visitor.controller");
const views = require("./Views.controller");
const github = require("./Github.controller");

function getTable(name) {
  switch (name) {
    case "Visitor":
      return visitor;

    case "Views":
      return views;

    case "Github":
      return github;
  }
}

// REDIRECT GET Request to the <Table>
exports.findAll = (req, res) => {
  let table = getTable(req.params.table);
  delete req.params.table;

  if (table) table.findAll(req, res);
  else this.logError || res.status(404).send({ message: "Table not Found" });
};

// REDIRECT GET Request to the <Table>
exports.findOne = (req, res) => {
  let table = getTable(req.params.table);
  delete req.params.table;

  if (table) table.findOne(req, res);
  else this.logError || res.status(404).send({ message: "Table not Found" });
};

// REDIRECT POST Request to the <Table>
exports.create = (req, res) => {
  let table = getTable(req.params.table);
  delete req.params.table;

  if (table) table.create(req, res);
  else this.logError || res.status(404).send({ message: "Table not Found" });
};

// REDIRECT PUT Request to the <Table>
exports.update = (req, res) => {
  let table = getTable(req.params.table);
  delete req.params.table;

  if (table) table.update(req, res);
  else this.logError || res.status(404).send({ message: "Table not Found" });
};

// REDIRECT DELETE Request to the <Table>
exports.destroy = (req, res) => {
  let table = getTable(req.params.table);
  delete req.params.table;

  if (table) table.destroy(req, res);
  else this.logError || res.status(404).send({ message: "Table not Found" });
};
