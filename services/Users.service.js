const { Users } = require("../models/index.js");

exports.find = (query, callback, err) => {
  Users.findAll({ where: query }).then(callback).catch(err);
};
