package config

import (
	models "api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Db *gorm.DB
}

func (o *Database) Init(env *Env) {
	var err error
	o.Db, err = gorm.Open(postgres.Open(
		"host="+env.DBHost+
			" user="+env.DBUser+
			" password="+env.DBPass+
			" port="+env.DBPort+
			" dbname="+env.DBName), &gorm.Config{})
	o.CheckOnErr(&err, "Failed on db connection")
}

func (o *Database) Migrate(bForce bool) {
	if bForce {
		o.Db.Migrator().DropTable(
			&models.Info{},
			&models.World{})
	}

	o.Db.AutoMigrate(
		&models.Info{},
		&models.World{})
}

func (*Database) CheckOnErr(err *error, msg string) {
	if *err != nil {
		panic(msg)
	}
}
