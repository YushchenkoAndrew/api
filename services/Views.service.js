const { Views } = require("../models/index.js");

exports.find = (query, callback, err) => {
  Views.findAll({ where: query }).then(callback).catch(err);
};

exports.create = (data, callback, err) => {
  Views.create(data).then(callback).catch(err);
};

exports.update = (data, query, callback, err) => {
  Views.update(data, { where: query }).then(callback).catch(err);
};

exports.destroy = (query, callback, err) => {
  Views.destroy({ where: query }).then(callback).catch(err);
};

exports.toString = () => "Views";

class Data {
  constructor({ Curr_Date, Count }) {
    this.Curr_Date = Curr_Date;
    this.Count = Count;
  }

  check() {
    return !this.Curr_Date && !this.Count;
  }
}

exports.Data = Data;
