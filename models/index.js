const config = require("../config/config");
const { Sequelize, Op } = require("sequelize");
const db = {};

// Set configurations
const sequelize = new Sequelize(config.database, config.username, config.password, {
  host: config.host,
  dialect: config.dialect,
  logging: config.logging,
});

sequelize
  .authenticate()
  .then((err) => console.log("\n\x1b[32m[INFO]\x1b[0m Connected to DataBase " + config.database))
  .catch((err) => console.log("\n\x1b[31m[INFO]\x1b[0m Unable to connect to db" + err));

db.Sequelize = Sequelize;
db.sequelize = sequelize;
db.Op = Op;

// Table models
db.Visitors = require("./Visitor.model.js")(sequelize, Sequelize);
db.Views = require("./Views.model.js")(sequelize, Sequelize);
db.Github = require("./Github.model.js")(sequelize, Sequelize);
db.UsedLanguages = require("./UsedLanguages.model.js")(sequelize, Sequelize);
db.Github.hasMany(db.UsedLanguages);

// Synchronize all Models
db.sequelize.sync().then(() => console.log("\x1b[32m[INFO]\x1b[0m Sync All defined DB Models"));

module.exports = db;
