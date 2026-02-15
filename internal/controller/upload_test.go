package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
)

func setupUploadController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&model.Video{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	logger := zerolog.Nop()
	shutdown := make(chan int, 1)
	ctrl := New(shutdown).(*Controller)
	ctrl.LoggerRegister(&logger)
	ctrl.DatabaseRegister(db)
	ctrl.ConfigurationRegister(&config.Config{})

	return ctrl
}

func TestController_UploadImport_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupUploadController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/upload/import", strings.NewReader(`{"path":"/media/triage/video.mp4","import_as":"v"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.UploadImport(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["id"] == nil || body["id"] == "" {
		t.Error("expected id in response")
	}
	if body["redirect"] == nil || body["redirect"] == "" {
		t.Error("expected redirect in response")
	}

	var vid model.Video
	if ctrl.datastore.First(&vid).RowsAffected < 1 {
		t.Error("video not created in DB")
	}
	if vid.Type != "v" {
		t.Errorf("expected type v, got %q", vid.Type)
	}
}

func TestController_UploadImport_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupUploadController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/upload/import", strings.NewReader(`{invalid`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.UploadImport(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
