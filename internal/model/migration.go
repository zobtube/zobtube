package model

import (
	"errors"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/logger"

	"gitlab.com/zobtube/zobtube/internal/config"
)

var modelToMigrate = []interface{}{
	Actor{},
	ActorAlias{},
	ActorLink{},
	Channel{},
	Video{},
	User{},
	UserSession{},
}

func New(cfg *config.Config) (db *gorm.DB, err error) {
	if cfg.DB.Driver == "sqlite" {
		db, err = gorm.Open(sqlite.Open(cfg.DB.Connstring), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else if cfg.DB.Driver == "postgresql" {
		db, err = gorm.Open(postgres.Open(cfg.DB.Connstring), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else {
		return db, errors.New("unsupported driver:" + cfg.DB.Driver)
	}

	if err != nil {
		return nil, err
	}

	// migrate all known models
	for _, m := range modelToMigrate {
		err = db.AutoMigrate(&m)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
