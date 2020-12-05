module.exports = (sequelize, Sequelize) => {
  return sequelize.define(
    "Users",
    {
      id: { type: Sequelize.DataTypes.INTEGER, primaryKey: true, autoIncrement: true },
      user: { type: Sequelize.DataTypes.STRING },
      pass: { type: Sequelize.DataTypes.STRING },
      role: { type: Sequelize.DataTypes.STRING },
    },
    {
      timestamps: false,
      freezeTableName: true,
    }
  );
};
