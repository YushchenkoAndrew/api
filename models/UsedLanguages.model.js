module.exports = (sequelize, Sequelize) => {
  return sequelize.define(
    "UsedLanguages",
    {
      id: { type: Sequelize.DataTypes.INTEGER, primaryKey: true, autoIncrement: true },
      Name: { type: Sequelize.DataTypes.STRING },
      Size: { type: Sequelize.DataTypes.INTEGER },
      GithubId: { type: Sequelize.DataTypes.INTEGER, references: { model: "Github", key: "id" } },
    },
    {
      timestamps: false,
      freezeTableName: true,
    }
  );
};
