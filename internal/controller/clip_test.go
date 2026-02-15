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

func setupClipController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&model.Video{}, &model.Actor{}, &model.CategorySub{}); err != nil {
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

func TestController_ClipView_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupClipController(t)
	vid := &model.Video{Name: "Clip1", Filename: "c1.mp4", Type: "c"}
	ctrl.datastore.Create(vid)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: vid.ID}}
	c.Request = httptest.NewRequest("GET", "/api/clip/"+vid.ID, nil)

	ctrl.ClipView(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["video"] == nil {
		t.Error("expected video in response")
	}
	if body["clip_ids"] == nil {
		t.Error("expected clip_ids in response")
	}
}

func TestController_ClipView_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupClipController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
	c.Request = httptest.NewRequest("GET", "/api/clip/nonexistent", nil)

	ctrl.ClipView(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestController_ClipView_NotAClip(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupClipController(t)
	vid := &model.Video{Name: "Video", Filename: "v.mp4", Type: "v"}
	ctrl.datastore.Create(vid)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: vid.ID}}
	c.Request = httptest.NewRequest("GET", "/api/clip/"+vid.ID, nil)

	ctrl.ClipView(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for non-clip type, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if errMsg, _ := body["error"].(string); errMsg != "not a clip" {
		t.Errorf("expected 'not a clip' error, got %v", errMsg)
	}
}
