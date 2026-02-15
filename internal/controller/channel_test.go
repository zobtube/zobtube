package controller

import (
	"encoding/json"
	"io"
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

func setupChannelController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&model.Channel{}, &model.Video{}); err != nil {
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

func TestController_ChannelList_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupChannelController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/channel", nil)

	ctrl.ChannelList(c)

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

func TestController_ChannelList_WithItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupChannelController(t)
	ctrl.datastore.Create(&model.Channel{Name: "Ch1"})
	ctrl.datastore.Create(&model.Channel{Name: "Ch2"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/channel", nil)

	ctrl.ChannelList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if total, _ := body["total"].(float64); total != 2 {
		t.Errorf("expected total 2, got %v", total)
	}
}

func TestController_ChannelGet_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupChannelController(t)
	ch := &model.Channel{Name: "My Channel"}
	ctrl.datastore.Create(ch)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: ch.ID}}
	c.Request = httptest.NewRequest("GET", "/api/channel/"+ch.ID, nil)

	ctrl.ChannelGet(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["channel"] == nil {
		t.Error("expected channel in response")
	}
}

func TestController_ChannelGet_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupChannelController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
	c.Request = httptest.NewRequest("GET", "/api/channel/nonexistent", nil)

	ctrl.ChannelGet(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestController_ChannelCreate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupChannelController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/channel", strings.NewReader(`{"name":"New Channel"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.ChannelCreate(c)

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
	if redir, _ := body["redirect"].(string); redir == "" {
		t.Error("expected redirect in response")
	}

	var ch model.Channel
	if ctrl.datastore.First(&ch).RowsAffected < 1 {
		t.Error("channel not created in DB")
	}
	if ch.Name != "New Channel" {
		t.Errorf("expected name 'New Channel', got %q", ch.Name)
	}
}

func TestController_ChannelCreate_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupChannelController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/channel", strings.NewReader(`{invalid`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.ChannelCreate(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestController_ChannelUpdate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupChannelController(t)
	ch := &model.Channel{Name: "Old Name"}
	ctrl.datastore.Create(ch)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: ch.ID}}
	c.Request = httptest.NewRequest("PUT", "/api/channel/"+ch.ID, strings.NewReader(`{"name":"New Name"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.ChannelUpdate(c)

	if w.Code != http.StatusOK {
		body, _ := io.ReadAll(w.Result().Body)
		t.Fatalf("expected 200, got %d: %s", w.Code, string(body))
	}
	var updated model.Channel
	ctrl.datastore.First(&updated, "id = ?", ch.ID)
	if updated.Name != "New Name" {
		t.Errorf("expected name 'New Name', got %q", updated.Name)
	}
}

func TestController_ChannelUpdate_PartialNilName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupChannelController(t)
	ch := &model.Channel{Name: "Keep Me"}
	ctrl.datastore.Create(ch)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: ch.ID}}
	c.Request = httptest.NewRequest("PUT", "/api/channel/"+ch.ID, strings.NewReader(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.ChannelUpdate(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var updated model.Channel
	ctrl.datastore.First(&updated, "id = ?", ch.ID)
	if updated.Name != "Keep Me" {
		t.Errorf("expected name unchanged 'Keep Me', got %q", updated.Name)
	}
}

func TestController_ChannelUpdate_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupChannelController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
	c.Request = httptest.NewRequest("PUT", "/api/channel/nonexistent", strings.NewReader(`{"name":"X"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.ChannelUpdate(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
