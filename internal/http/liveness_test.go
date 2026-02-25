package http_test

import (
	"embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/controller"
	httpserver "github.com/zobtube/zobtube/internal/http"
	"github.com/zobtube/zobtube/internal/model"
)

//go:embed liveness_test.go
var embedFS embed.FS

func TestPingRoute(t *testing.T) {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	var c chan int
	cont := controller.New(c)
	server := httpserver.New(&embedFS, false, &logger)
	server.ControllerSetupDefault(&cont)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "alive", w.Body.String())
}

func TestFailsafeUnexpectedErrorPingRoute(t *testing.T) {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	var c chan int
	cont := controller.New(c)
	server := httpserver.New(&embedFS, false, &logger)
	server.ControllerSetupFailsafeError(cont, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "ko", w.Body.String())
}

func TestUserIsAuthenticated_NoCredentials_Returns401(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.UserSession{}, &model.ApiToken{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	shutdown := make(chan int, 1)
	ctrl := controller.New(shutdown).(*controller.Controller)
	ctrl.DatabaseRegister(db)
	ctrl.ConfigurationRegister(&config.Config{Authentication: true})
	zlog := zerolog.Nop()
	ctrl.LoggerRegister(&zlog)

	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	server := httpserver.New(&embedFS, false, &logger)
	var iface controller.AbstractController = ctrl
	server.ControllerSetupDefault(&iface)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/profile", nil)
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	assert.Equal(t, "unauthorized", body["error"])
}
