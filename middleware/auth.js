const jwt = require("jsonwebtoken");
const { logRequest, logInfo } = require("../lib/log.js");
const { errorHandler, resByType, getDataByType } = require("../lib/resHandler");
const Users = require("../services/Users.service.js");
const md5 = require("./md5");

exports.authorizationToken = (req, res, next) => {
  logRequest("GET", "Authorization ...");
  let token = req.header("Authorization")?.split(" ");

  if (process.env.API_DEBUG) return logInfo("DEBUG MODE") || next();
  if (!token || !token[1]) return errorHandler(401, "Wrong Token declaration", req, res);

  jwt.verify(token[1], process.env.SECRET_TOKEN, (err, user) => {
    if (err) return errorHandler(403, err.message, req, res);

    req.user = user;
    logInfo("SUCCESS", { user: user.user, role: user.role });
    next();
  });
};

exports.generateToken = (req, res) => {
  logRequest("POST", "GENERATE TOKEN; DATA = ", req.body);
  const data = getDataByType(req.header("Content-Type"), req.body);

  Users.find(
    // Query
    null,
    // Callback
    (res) => {
      let found = false;
      for (let i in res) {
        if (md5(res[i].pass + (data?.rand ?? "")) == data?.pass) {
          found = true;
          break;
        }
      }

      if (!found) return errorHandler(406, "Incorrect User or Pass Value", req, res);
      const accessToken = jwt.sign(data[0].dataValues, process.env.SECRET_TOKEN, { expiresIn: process.env.TOKEN_EXPIRE || "600s" });
      resByType(req.header("Content-Type"), { accessToken }, res);
    },
    // Error
    (err) => errorHandler(500, err.message, req, res)
  );
};
