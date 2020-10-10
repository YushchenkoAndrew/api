const { Github } = require("../models/index.js");

exports.find = (query, callback, err) => {
  Github.findAll({ where: query }).then(callback).catch(err);
};

exports.create = (data, callback, err) => {
  Github.create(data).then(callback).catch(err);
};

exports.update = (data, query, callback, err) => {
  Github.update(data, { where: query }).then(callback).catch(err);
};

exports.destroy = (query, callback, err) => {
  Github.destroy({ where: query }).then(callback).catch(err);
};

exports.toString = () => "Github";

class Data {
  constructor({ Name, UpdateAt, Context }) {
    this.Name = Name;
    this.UpdateAt = UpdateAt;
    this.Context = Context;
  }

  check() {
    return this.Name === undefined || this.UpdateAt === undefined || this.Context === undefined;
  }
}

exports.Data = Data;
