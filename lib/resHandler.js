const jsonToXml = require("jsontoxml");
const { logError, logDebug } = require("./log");

exports.errorHandler = (status, message, req, res) => {
  logError(status, message);
  if (res) this.resByType(req.header("Content-Type"), { message: message }, res, status);
};

exports.resByType = (type, data, res, status = 200) => res.status(status).send(type && type.endsWith("xml") ? jsonToXml(JSON.stringify(data)) : data);
