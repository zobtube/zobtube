package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
)

func TestMigrationSQLite(t *testing.T) {
	// prepare configuration for test database
	cfg := config.Config{}
	cfg.DB.Driver = "sqlite"
	cfg.DB.Connstring = "./zt-test.sqlite"

	_, err := model.New(&cfg)

	assert.Equal(t, nil, err)
}
