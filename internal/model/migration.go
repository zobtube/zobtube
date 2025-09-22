package model

import (
	"errors"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zobtube/zobtube/internal/config"
)

var modelToMigrate = []any{
	Actor{},
	ActorAlias{},
	ActorLink{},
	Category{},
	CategorySub{},
	Channel{},
	Configuration{},
	Video{},
	VideoView{},
	Task{},
	User{},
	UserSession{},
}

func New(cfg *config.Config) (db *gorm.DB, err error) {
	switch cfg.DB.Driver {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.DB.Connstring), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	case "postgresql":
		db, err = gorm.Open(postgres.Open(cfg.DB.Connstring), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	default:
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
