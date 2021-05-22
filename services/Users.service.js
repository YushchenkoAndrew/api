const { Users } = require("../models/index.js");

exports.find = (query, callback, err) => {
  Users.findAll(query ? { where: query } : {})
    .then(callback)
    .catch(err);
};
