package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zobtube/zobtube/internal/config"
)

func isSQLite(cfg *config.Config) bool {
	return strings.EqualFold(cfg.DB.Driver, "sqlite")
}

var modelToMigrate = []any{
	Actor{},
	ActorAlias{},
	ActorLink{},
	Category{},
	CategorySub{},
	Channel{},
	Configuration{},
	Library{},
	Provider{},
	Video{},
	VideoView{},
	Task{},
	User{},
	UserSession{},
	ApiToken{},
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

	if isSQLite(cfg) {
		if err := migrateSQLiteVideoLibraryID(db); err != nil {
			return nil, err
		}
	}

	// migrate all known models
	for _, m := range modelToMigrate {
		err = db.AutoMigrate(&m)
		if err != nil {
			log.Error().Err(err).Msg("failed to migrate model")
			fmt.Printf("failed to migrate model: model=%T error=%v\n", m, err)
			if isSQLite(cfg) && isSQLiteNotNullDefaultNullError(err) {
				// GORM tried to add library_id without default; add it with default then retry
				if retryErr := migrateSQLiteVideoLibraryID(db); retryErr != nil {
					return nil, retryErr
				}
				err = db.AutoMigrate(&m)
			}
			if err != nil {
				return nil, err
			}
		}
	}

	return db, nil
}

// migrateSQLiteVideoLibraryID adds videos.library_id with DEFAULT for SQLite so that
// adding the column to a non-empty table succeeds. Run before AutoMigrate so
// GORM does not try to add the column without a default.
func migrateSQLiteVideoLibraryID(db *gorm.DB) error {
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='videos'").Scan(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	const defaultLibraryUUID = "00000000-0000-0000-0000-000000000000"
	// Use literal default so SQLite clearly has a non-null default when adding the column.
	err := db.Exec("ALTER TABLE videos ADD COLUMN library_id TEXT DEFAULT '" + defaultLibraryUUID + "'").Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate column name") {
			return nil
		}
		return err
	}
	_ = db.Exec("CREATE INDEX IF NOT EXISTS idx_videos_library_id ON videos(library_id)").Error
	return nil
}

func isSQLiteNotNullDefaultNullError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Cannot add a NOT NULL column with default value NULL")
}
