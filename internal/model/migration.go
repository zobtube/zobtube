package model

import (
	"errors"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gitlab.com/zobtube/zobtube/internal/config"
)

var modelToMigrate = []interface{}{
	Actor{},
	ActorAlias{},
	ActorLink{},
	Channel{},
	Video{},
}

func New(cfg *config.Config) (db *gorm.DB, err error) {
	if cfg.DbDriver == "sqlite" {
		db, err = gorm.Open(sqlite.Open(cfg.DbConnstring), &gorm.Config{})
	} else if cfg.DbDriver == "postgresql" {
		db, err = gorm.Open(postgres.Open(cfg.DbConnstring), &gorm.Config{})
	} else {
		return db, errors.New("unsupported driver:" + cfg.DbDriver)
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
