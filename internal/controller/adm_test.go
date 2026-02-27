package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/runner"
	"github.com/zobtube/zobtube/internal/storage"
	"github.com/zobtube/zobtube/internal/task/common"
)

func setupAdmController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(
		&model.Video{}, &model.Actor{}, &model.Channel{}, &model.Category{},
		&model.User{}, &model.Task{}, &model.Provider{}, &model.Configuration{},
		&model.ApiToken{}, &model.Library{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	defaultLibID, err := model.EnsureDefaultLibrary(db, "/tmp")
	if err != nil {
		t.Fatalf("failed to ensure default library: %v", err)
	}
	logger := zerolog.Nop()
	shutdown := make(chan int, 1)
	ctrl := New(shutdown).(*Controller)
	ctrl.LoggerRegister(&logger)
	ctrl.DatabaseRegister(db)
	cfg := &config.Config{}
	cfg.DB.Driver = "sqlite"
	cfg.DB.Connstring = ":memory:"
	cfg.DefaultLibraryID = defaultLibID
	ctrl.ConfigurationRegister(cfg)
	ctrl.BuildDetailsRegister("0.0.0", "abc123", "2024-01-01")
	storageResolver := storage.NewResolver(db)
	ctrl.StorageResolverRegister(storageResolver)

	// Runner for AdmTaskRetry
	r := &runner.Runner{}
	r.RegisterTask(&common.Task{
		Name:  "adm-test-task",
		Steps: []common.Step{{Name: "noop", NiceName: "No-op", Func: func(*common.Context, common.Parameters) (string, error) { return "", nil }}},
	})
	r.Start(cfg, db, storageResolver)
	ctrl.RunnerRegister(r)

	return ctrl
}

func TestController_AdmHome(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	ctrl.datastore.Create(&model.Video{Name: "V", Filename: "v.mp4", Type: "v"})
	ctrl.datastore.Create(&model.Actor{Name: "A"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/adm", nil)

	ctrl.AdmHome(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["video_count"] != float64(1) {
		t.Errorf("expected video_count 1, got %v", body["video_count"])
	}
	if body["actor_count"] != float64(1) {
		t.Errorf("expected actor_count 1, got %v", body["actor_count"])
	}
}

func TestController_AdmUserNew_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/adm/user", strings.NewReader(`{"username":"newuser","password":"secret","admin":false}`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.AdmUserNew(c)

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
}

func TestController_AdmUserNew_EmptyPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/adm/user", strings.NewReader(`{"username":"u","password":"","admin":false}`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.AdmUserNew(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if errMsg, _ := body["error"].(string); errMsg != "password cannot be empty" {
		t.Errorf("expected password error, got %v", errMsg)
	}
}

func TestController_AdmUserNew_DuplicateUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	ctrl.datastore.Create(&model.User{Username: "taken"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/adm/user", strings.NewReader(`{"username":"taken","password":"secret","admin":false}`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.AdmUserNew(c)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if errMsg, _ := body["error"].(string); errMsg != "username already taken" {
		t.Errorf("expected username taken error, got %v", errMsg)
	}
}

func TestController_AdmTaskRetry_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	task := &model.Task{Name: "adm-test-task", Status: model.TaskStatusError}
	ctrl.datastore.Create(task)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: task.ID}}
	c.Request = httptest.NewRequest("POST", "/api/adm/task/"+task.ID+"/retry", nil)

	ctrl.AdmTaskRetry(c)

	if w.Code != http.StatusOK {
		body, _ := io.ReadAll(w.Result().Body)
		t.Fatalf("expected 200, got %d: %s", w.Code, string(body))
	}
}

func TestController_AdmTaskRetry_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
	c.Request = httptest.NewRequest("POST", "/api/adm/task/nonexistent/retry", nil)

	ctrl.AdmTaskRetry(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestController_AdmTaskView_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	task := &model.Task{Name: "T", Status: model.TaskStatusTodo}
	ctrl.datastore.Create(task)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: task.ID}}
	c.Request = httptest.NewRequest("GET", "/api/adm/task/"+task.ID, nil)

	ctrl.AdmTaskView(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestController_AdmTaskView_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
	c.Request = httptest.NewRequest("GET", "/api/adm/task/nonexistent", nil)

	ctrl.AdmTaskView(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestController_AdmUserDelete_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}
	c.Request = httptest.NewRequest("DELETE", "/api/adm/user/nonexistent", nil)

	ctrl.AdmUserDelete(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestController_AdmVideoList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	ctrl.datastore.Create(&model.Video{Name: "V", Filename: "v.mp4", Type: "v"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/adm/video", nil)

	ctrl.AdmVideoList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if total, _ := body["total"].(float64); total != 1 {
		t.Errorf("expected total 1, got %v", total)
	}
}

func TestController_AdmCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/adm/category", nil)

	ctrl.AdmCategory(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestController_AdmTokenList_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/adm/tokens", nil)

	ctrl.AdmTokenList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	tokens, _ := body["tokens"].([]any)
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(tokens))
	}
}

func TestController_AdmTokenList_WithTokens(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	u1 := &model.User{Username: "user1"}
	u2 := &model.User{Username: "user2"}
	ctrl.datastore.Create(u1)
	ctrl.datastore.Create(u2)
	ctrl.datastore.Create(&model.ApiToken{UserID: u1.ID, Name: "Token A", TokenHash: "h1", CreatedAt: time.Now()})
	ctrl.datastore.Create(&model.ApiToken{UserID: u2.ID, Name: "Token B", TokenHash: "h2", CreatedAt: time.Now()})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/adm/tokens", nil)

	ctrl.AdmTokenList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	tokens, _ := body["tokens"].([]any)
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	for i, raw := range tokens {
		tok, _ := raw.(map[string]any)
		if tok["id"] == nil || tok["name"] == nil || tok["username"] == nil {
			t.Errorf("token %d: missing id, name, or username", i)
		}
		if tok["token_hash"] != nil {
			t.Error("response must not contain token_hash")
		}
	}
}

func TestController_AdmTokenDelete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)
	tok := &model.ApiToken{UserID: user.ID, Name: "T", TokenHash: "h", CreatedAt: time.Now()}
	ctrl.datastore.Create(tok)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: tok.ID}}
	c.Request = httptest.NewRequest("DELETE", "/api/adm/tokens/"+tok.ID, nil)

	ctrl.AdmTokenDelete(c)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	var count int64
	ctrl.datastore.Model(&model.ApiToken{}).Where("id = ?", tok.ID).Count(&count)
	if count != 0 {
		t.Error("token should be deleted")
	}
}

func TestController_AdmTokenDelete_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAdmController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "00000000-0000-0000-0000-000000000000"}}
	c.Request = httptest.NewRequest("DELETE", "/api/adm/tokens/00000000-0000-0000-0000-000000000000", nil)

	ctrl.AdmTokenDelete(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
