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

func setupVideoViewController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&model.Video{}, &model.VideoView{}, &model.User{}); err != nil {
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

func TestController_VideoViewIncrement_FirstView(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoViewController(t)
	vid := &model.Video{Name: "V", Filename: "v.mp4", Type: "v"}
	user := &model.User{Username: "viewer"}
	ctrl.datastore.Create(vid)
	ctrl.datastore.Create(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: vid.ID}}
	c.Request = httptest.NewRequest("POST", "/api/video/"+vid.ID+"/view", nil)
	c.Set("user", user)

	ctrl.VideoViewIncrement(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["view-count"] != float64(1) {
		t.Errorf("expected view-count 1, got %v", body["view-count"])
	}
}

func TestController_VideoViewIncrement_SubsequentView(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoViewController(t)
	vid := &model.Video{Name: "V", Filename: "v.mp4", Type: "v"}
	user := &model.User{Username: "viewer"}
	ctrl.datastore.Create(vid)
	ctrl.datastore.Create(user)
	ctrl.datastore.Create(&model.VideoView{VideoID: vid.ID, UserID: user.ID, Count: 3})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: vid.ID}}
	c.Request = httptest.NewRequest("POST", "/api/video/"+vid.ID+"/view", nil)
	c.Set("user", user)

	ctrl.VideoViewIncrement(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["view-count"] != float64(4) {
		t.Errorf("expected view-count 4, got %v", body["view-count"])
	}
}

func TestController_VideoViewIncrement_VideoNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupVideoViewController(t)
	user := &model.User{Username: "viewer"}
	ctrl.datastore.Create(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
	c.Request = httptest.NewRequest("POST", "/api/video/nonexistent/view", nil)
	c.Set("user", user)

	ctrl.VideoViewIncrement(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
