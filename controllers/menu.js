const { logRequest, logDebug } = require("../lib/log");
const { resByType } = require("../lib/resHandler");
const config = require("../config/config.menu.js");
const db = require("../models/index");

exports.menu = (req, res) => {
  logRequest("GET", "API Navigation");
  let url = `${req.hostname}:${process.env.API_PORT}`;

  let data = {};
  for (let key in config) data[key] = url + config[key];

  resByType(req.header("Content-Type"), data, res);
};

exports.getTables = (req, res) => {
  logRequest("GET", "LIST Tables");

  res.send(db.Tables);
};
