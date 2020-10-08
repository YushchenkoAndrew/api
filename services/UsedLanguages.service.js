const { UsedLanguages } = require("../models/index.js");

exports.find = (query, callback, err) => {
  UsedLanguages.findAll({ where: query }).then(callback).catch(err);
};

exports.create = (data, callback, err) => {
  UsedLanguages.create(data).then(callback).catch(err);
};

exports.update = (data, query, callback, err) => {
  UsedLanguages.update(data, { where: query }).then(callback).catch(err);
};

exports.destroy = (query, callback, err) => {
  UsedLanguages.destroy({ where: query }).then(callback).catch(err);
};

exports.toString = () => "UsedLanguages";

class Data {
  constructor({ Name, Size, RepoID }) {
    this.Name = Name;
    this.Size = Size;
    this.RepoID = RepoID;
  }

  check() {
    return !this.Name && !this.Size && !this.RepoID;
  }
}

exports.Data = Data;
