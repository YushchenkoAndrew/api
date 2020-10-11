const { logError } = require("./log");

exports.errorHandler = (status, message, req, res) => {
  logError(status, message);
  if (res) res.status(status).send({ message: message });
};
