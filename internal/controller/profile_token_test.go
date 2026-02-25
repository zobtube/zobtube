package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
)

func setupProfileTokenController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.ApiToken{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	logger := zerolog.Nop()
	shutdown := make(chan int, 1)
	ctrl := New(shutdown).(*Controller)
	ctrl.LoggerRegister(&logger)
	ctrl.DatabaseRegister(db)
	ctrl.ConfigurationRegister(&config.Config{Authentication: true})

	return ctrl
}

func TestController_ResolveUserByApiTokenHash_ValidToken(t *testing.T) {
	ctrl := setupProfileTokenController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)
	rawToken := "test-token-secret"
	sum := sha256.Sum256([]byte(rawToken))
	hash := hex.EncodeToString(sum[:])
	ctrl.datastore.Create(&model.ApiToken{UserID: user.ID, Name: "t1", TokenHash: hash, CreatedAt: time.Now()})

	resolved, ok := ctrl.ResolveUserByApiTokenHash(rawToken)
	if !ok || resolved == nil {
		t.Fatalf("expected user, got ok=%v user=%v", ok, resolved)
	}
	if resolved.ID != user.ID || resolved.Username != "u" {
		t.Errorf("expected user %s/%s, got %s/%s", user.ID, user.Username, resolved.ID, resolved.Username)
	}
}

func TestController_ResolveUserByApiTokenHash_InvalidToken(t *testing.T) {
	ctrl := setupProfileTokenController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)
	sum := sha256.Sum256([]byte("real-token"))
	ctrl.datastore.Create(&model.ApiToken{UserID: user.ID, Name: "t1", TokenHash: hex.EncodeToString(sum[:]), CreatedAt: time.Now()})

	resolved, ok := ctrl.ResolveUserByApiTokenHash("wrong-token")
	if ok || resolved != nil {
		t.Errorf("expected no user, got ok=%v user=%v", ok, resolved)
	}
}

func TestController_ResolveUserByApiTokenHash_EmptyToken(t *testing.T) {
	ctrl := setupProfileTokenController(t)
	resolved, ok := ctrl.ResolveUserByApiTokenHash("")
	if ok || resolved != nil {
		t.Errorf("expected no user for empty token, got ok=%v user=%v", ok, resolved)
	}
}

func TestController_ProfileTokenList_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileTokenController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/profile/tokens", nil)
	c.Set("user", user)

	ctrl.ProfileTokenList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	tokens, _ := resp["tokens"].([]any)
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(tokens))
	}
}

func TestController_ProfileTokenList_WithTokens(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileTokenController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)
	ctrl.datastore.Create(&model.ApiToken{UserID: user.ID, Name: "A", TokenHash: "h1", CreatedAt: time.Now()})
	ctrl.datastore.Create(&model.ApiToken{UserID: user.ID, Name: "B", TokenHash: "h2", CreatedAt: time.Now()})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/profile/tokens", nil)
	c.Set("user", user)

	ctrl.ProfileTokenList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	tokens, _ := resp["tokens"].([]any)
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	// Response must not contain token or hash
	bodyStr := w.Body.String()
	if strings.Contains(bodyStr, "h1") || strings.Contains(bodyStr, "h2") {
		t.Error("response must not contain token hash")
	}
}

func TestController_ProfileTokenCreate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileTokenController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/profile/tokens", strings.NewReader(`{"name":"My token"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user", user)

	ctrl.ProfileTokenCreate(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["id"] == nil || resp["id"] == "" {
		t.Error("expected id in response")
	}
	if resp["name"] != "My token" {
		t.Errorf("expected name My token, got %v", resp["name"])
	}
	token, _ := resp["token"].(string)
	if token == "" || len(token) != 64 {
		t.Errorf("expected 64-char hex token, got %q", token)
	}
	// Resolver should find user with this token
	resolved, ok := ctrl.ResolveUserByApiTokenHash(token)
	if !ok || resolved.ID != user.ID {
		t.Errorf("ResolveUserByApiTokenHash with created token: got ok=%v user=%v", ok, resolved)
	}
}

func TestController_ProfileTokenCreate_EmptyName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileTokenController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/profile/tokens", strings.NewReader(`{"name":""}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user", user)

	ctrl.ProfileTokenCreate(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if errMsg, _ := resp["error"].(string); errMsg != "name is required" {
		t.Errorf("expected error 'name is required', got %q", errMsg)
	}
}

func TestController_ProfileTokenCreate_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileTokenController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/profile/tokens", strings.NewReader(`{invalid`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user", user)

	ctrl.ProfileTokenCreate(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestController_ProfileTokenDelete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileTokenController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)
	tok := &model.ApiToken{UserID: user.ID, Name: "T", TokenHash: "h", CreatedAt: time.Now()}
	ctrl.datastore.Create(tok)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/api/profile/tokens/"+tok.ID, nil)
	c.Params = gin.Params{{Key: "id", Value: tok.ID}}
	c.Set("user", user)

	ctrl.ProfileTokenDelete(c)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
	var count int64
	ctrl.datastore.Model(&model.ApiToken{}).Where("id = ?", tok.ID).Count(&count)
	if count != 0 {
		t.Error("token should be deleted")
	}
}

func TestController_ProfileTokenDelete_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileTokenController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/api/profile/tokens/00000000-0000-0000-0000-000000000000", nil)
	c.Params = gin.Params{{Key: "id", Value: "00000000-0000-0000-0000-000000000000"}}
	c.Set("user", user)

	ctrl.ProfileTokenDelete(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestController_ProfileTokenDelete_WrongUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileTokenController(t)
	user1 := &model.User{Username: "u1"}
	user2 := &model.User{Username: "u2"}
	ctrl.datastore.Create(user1)
	ctrl.datastore.Create(user2)
	tok := &model.ApiToken{UserID: user1.ID, Name: "T", TokenHash: "h", CreatedAt: time.Now()}
	ctrl.datastore.Create(tok)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/api/profile/tokens/"+tok.ID, nil)
	c.Params = gin.Params{{Key: "id", Value: tok.ID}}
	c.Set("user", user2) // user2 trying to delete user1's token

	ctrl.ProfileTokenDelete(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 (forbidden as not found for other user), got %d: %s", w.Code, w.Body.String())
	}
	var count int64
	ctrl.datastore.Model(&model.ApiToken{}).Where("id = ?", tok.ID).Count(&count)
	if count != 1 {
		t.Error("token must not be deleted")
	}
}
