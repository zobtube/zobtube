package http_test

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/zobtube/zobtube/internal/controller"
	httpserver "github.com/zobtube/zobtube/internal/http"
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
