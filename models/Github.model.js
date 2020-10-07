module.exports = (sequelize, Sequelize) => {
  return sequelize.define(
    "Github",
    {
      id: { type: Sequelize.DataTypes.INTEGER, primaryKey: true, autoIncrement: true },
      Name: { type: Sequelize.DataTypes.STRING },
      UpdateAt: { type: Sequelize.DataTypes.DATEONLY },
      Context: { type: Sequelize.DataTypes.STRING },
    },
    {
      timestamps: false,
      freezeTableName: true,
    }
  );
};
