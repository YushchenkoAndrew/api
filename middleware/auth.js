require("dotenv").config();
const jwt = require("jsonwebtoken");

exports.authorizationToken = (req, res, next) => {
  let token = req.header("Authorization");

  console.log("START");
  console.log(token);

  if (!token) return res.sendStatus(401);

  console.log("Authentication!!!");

  jwt.verify(token, process.env.SECRET_TOKEN, (err, user) => {
    console.log(err);

    if (err) return res.sendStatus(403);

    console.log(user);
    next();
  });
};

exports.generateToken = (req, res) => {
  res.send({ accessToken: jwt.sign({ user: "admin" }, process.env.SECRET_TOKEN, { expiresIn: process.env.TOKEN_EXPIRE || 600 }) });
};
