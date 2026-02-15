package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func setupErrController(t *testing.T) *Controller {
	logger := zerolog.Nop()
	shutdown := make(chan int, 1)
	ctrl := New(shutdown).(*Controller)
	ctrl.LoggerRegister(&logger)

	return ctrl
}

func TestController_ErrNotFound_Returns404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupErrController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/unknown", nil)

	ctrl.ErrNotFound(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["error"] != "not found" {
		t.Errorf("expected error 'not found', got %v", body["error"])
	}
}

func TestController_ErrUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupErrController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/restricted", nil)

	ctrl.ErrUnauthorized(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["error"] != "unauthorized" {
		t.Errorf("expected error 'unauthorized', got %v", body["error"])
	}
}
