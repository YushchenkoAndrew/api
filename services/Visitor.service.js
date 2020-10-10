const { Visitors } = require("../models/index.js");

exports.find = (query, callback, err) => {
  Visitors.findAll({ where: query }).then(callback).catch(err);
};

exports.create = (data, callback, err) => {
  Visitors.create(data).then(callback).catch(err);
};

exports.update = (data, query, callback, err) => {
  Visitors.update(data, { where: query }).then(callback).catch(err);
};

exports.destroy = (query, callback, err) => {
  Visitors.destroy({ where: query }).then(callback).catch(err);
};

exports.toString = () => "Visitor";

class Data {
  constructor({ Country, ip, Visit_Date, Count }) {
    this.Country = Country;
    this.ip = ip;
    this.Visit_Date = Visit_Date;
    this.Count = Count;
  }

  check() {
    return this.Country === undefined && this.ip === undefined && this.Visit_Date === undefined && this.Count === undefined;
  }
}

exports.Data = Data;
