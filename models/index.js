const { logInfo } = require("../lib/log");
const { errorHandler } = require("../lib/resHandler");
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
  .then((data) => logInfo("Connected to DataBase " + config.database))
  .catch((err) => errorHandler(500, `Unable to connect to db ${err.message}`));

db.Sequelize = Sequelize;
db.sequelize = sequelize;
db.Op = Op;

// Table models
db.Visitors = require("./Visitor.model.js")(sequelize, Sequelize);
db.Views = require("./Views.model.js")(sequelize, Sequelize);
db.Github = require("./Github.model.js")(sequelize, Sequelize);
db.UsedLanguages = require("./UsedLanguages.model.js")(sequelize, Sequelize);
db.Users = require("./Users.model.js")(sequelize, Sequelize);

db.Tables = ["Visitor", "Views", "Github", "UsedLangues"];

// Relation between databases
db.Github.hasMany(db.UsedLanguages);

// Synchronize all Models
db.sequelize.sync().then(() => logInfo("Sync All defined DB Models"));

module.exports = db;
