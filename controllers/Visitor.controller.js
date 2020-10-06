const { Visitors } = require("../models/index.js");

exports.test = (req, res) => {
  console.log("Testing... It's Working!!!");
  res.send({ message: "OK" });
};

exports.findAll = (req, res) => {
  let condition = Object.keys(req.query).length ? req.query : null;
  console.log("\x1b[34m[REQUEST]\x1b[0m findAll with condition =", condition);

  Visitors.findAll({ where: condition })
    .then((data) => res.send(data))
    .catch((err) => {
      console.log(`\x1b[31m[Error]\x1b[0m ${err}`);
      res.status(500).send({ message: err.message });
    });
};

exports.findOne = (req, res) => {
  Visitors.findAll({ where: { id: req.params.id } }).then((data) => res.send(data));
};

// exports.create = (table, data) => {
//   let newVisitor = {};

//   for (let i in data) {
//     let params = data[i].split("=");
//     newVisitor[params[0]] = params[1] || null;
//   }

//   return getTable(table).create(newVisitor);
// };

// exports.delete = (table, key, value) => {
//   let condition = {};
//   condition[key] = value;

//   return getTable(table).destroy({ where: condition });
// };

// exports.update = (table, data) => {
//   let condition = data.splice(0, 1)[0].split("=");
//   let where = {};
//   where[condition[0]] = condition[1];

//   let newValue = {};

//   for (let i in data) {
//     let params = data[i].split("=");
//     newValue[params[0]] = params[1];
//   }

//   return getTable(table).update(newValue, { where: where });
// };
