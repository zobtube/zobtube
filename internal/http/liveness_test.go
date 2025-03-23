package http_test

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zobtube/zobtube/internal/controller"
	httpServer "github.com/zobtube/zobtube/internal/http"
)

//go:embed liveness_test.go
var embedFS embed.FS

func TestPingRoute(t *testing.T) {
	var c chan int
	cont := controller.New(c)
	server, _ := httpServer.New(&cont, &embedFS)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	server.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "alive", w.Body.String())
}
