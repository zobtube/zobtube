package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
)

func setupHomeController(t *testing.T) *Controller {
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

func TestController_Home_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupHomeController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/home", nil)

	ctrl.Home(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if total, _ := body["total"].(float64); total != 0 {
		t.Errorf("expected total 0, got %v", total)
	}
}

func TestController_Home_WithVideos(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupHomeController(t)
	ctrl.datastore.Create(&model.Video{Name: "V1", Filename: "v1.mp4", Type: "v"})
	ctrl.datastore.Create(&model.Video{Name: "V2", Filename: "v2.mp4", Type: "v"})
	ctrl.datastore.Create(&model.Video{Name: "C1", Filename: "c1.mp4", Type: "c"}) // clip, not video

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/home", nil)

	ctrl.Home(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if total, _ := body["total"].(float64); total != 2 {
		t.Errorf("expected total 2 (videos only), got %v", total)
	}
}
