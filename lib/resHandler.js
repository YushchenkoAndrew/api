const jsonToXml = require("jsontoxml");
const { logError, logDebug } = require("./log");

exports.errorHandler = (status, message, req, res) => {
  logError(status, message);
  if (res) this.resByType(req.header("Content-Type"), { message: message }, res, status);
};

exports.resByType = (type, data, res, status = 200) => res.status(status).send(type?.endsWith("xml") ? jsonToXml(JSON.stringify(data)) : data);

exports.getDataByType = (type, data) => (type?.endsWith("xml") ? clearXmlData(data.root) : data);

function clearXmlData(obj) {
  if (Array.isArray(obj) && obj.length <= 1) return obj.join("");

  for (let key in obj || []) obj[key] = clearXmlData(obj[key]);
  return obj;
}
