package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zobtube/zobtube/internal/config"
)

func TestConfigFromFile(t *testing.T) {
	// create tmp configuration
	f, err := os.CreateTemp("", "config")
	assert.Equal(t, nil, err)
	defer os.Remove(f.Name())

	_, _ = f.WriteString("server:\n")
	_, _ = f.WriteString("  bind: 0.0.0.0:8080\n")
	_, _ = f.WriteString("media:\n")
	_, _ = f.WriteString("  path: library_test\n")
	_, _ = f.WriteString("db:\n")
	_, _ = f.WriteString("  driver: sqlite\n")
	_, _ = f.WriteString("  connstring: ./zt-test.sqlite\n")
	f.Close()

	// load configuration
	cfg, err := config.New(f.Name())
	assert.Equal(t, nil, err)
	assert.Equal(t, "0.0.0.0:8080", cfg.Server.Bind)
	assert.Equal(t, "library_test", cfg.Media.Path)
	assert.Equal(t, "sqlite", cfg.DB.Driver)
	assert.Equal(t, "./zt-test.sqlite", cfg.DB.Connstring)
}

func TestConfigFromEnv(t *testing.T) {
	// set config from env vars
	os.Setenv("ZT_SERVER_BIND", "0.0.0.0:8080")
	os.Setenv("ZT_MEDIA_PATH", "library_test")
	os.Setenv("ZT_DB_DRIVER", "sqlite")
	os.Setenv("ZT_DB_CONNSTRING", "./zt-test.sqlite")

	// load configuration
	cfg, err := config.New("null")
	assert.Equal(t, nil, err)
	assert.Equal(t, "0.0.0.0:8080", cfg.Server.Bind)
	assert.Equal(t, "library_test", cfg.Media.Path)
	assert.Equal(t, "sqlite", cfg.DB.Driver)
	assert.Equal(t, "./zt-test.sqlite", cfg.DB.Connstring)
}
