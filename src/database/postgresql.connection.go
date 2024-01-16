package database

import (
	"github.com/isd-sgcu/johnjud-backend/src/app/model/like"
	"github.com/isd-sgcu/johnjud-backend/src/app/model/pet"
	"github.com/isd-sgcu/johnjud-backend/src/app/model/user"
	"github.com/isd-sgcu/johnjud-backend/src/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func InitPostgresDatabase(conf *config.Database, isDebug bool) (db *gorm.DB, err error) {
	gormConf := &gorm.Config{}

	if !isDebug {
		gormConf.Logger = gormLogger.Default.LogMode(gormLogger.Silent)
	}

	db, err = gorm.Open(postgres.Open(conf.Url), gormConf)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&user.User{}, &like.Like{}, &pet.Pet{})
	if err != nil {
		return nil, err
	}

	return
}
