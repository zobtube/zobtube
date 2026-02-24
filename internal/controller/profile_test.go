package controller

import (
	"crypto/sha256"
	"encoding/hex"
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

func setupProfileController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
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

func TestController_ProfileChangePassword_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileController(t)
	currentSum := sha256.Sum256([]byte("oldpass"))
	user := &model.User{Username: "u", Password: hex.EncodeToString(currentSum[:])}
	ctrl.datastore.Create(user)

	body := `{"current_password":"oldpass","new_password":"newpass"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/profile/password", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user", user)

	ctrl.ProfileChangePassword(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["ok"] != true {
		t.Errorf("expected ok true, got %v", resp["ok"])
	}
	// Password in DB should be updated
	ctrl.datastore.First(user)
	newSum := sha256.Sum256([]byte("newpass"))
	expected := hex.EncodeToString(newSum[:])
	if user.Password != expected {
		t.Errorf("expected password hash %s, got %s", expected, user.Password)
	}
}

func TestController_ProfileChangePassword_WrongCurrentPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileController(t)
	currentSum := sha256.Sum256([]byte("realpass"))
	user := &model.User{Username: "u", Password: hex.EncodeToString(currentSum[:])}
	ctrl.datastore.Create(user)

	body := `{"current_password":"wrongpass","new_password":"newpass"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/profile/password", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user", user)

	ctrl.ProfileChangePassword(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if errMsg, _ := resp["error"].(string); errMsg != "wrong current password" {
		t.Errorf("expected error 'wrong current password', got %q", errMsg)
	}
	// Password in DB must be unchanged
	ctrl.datastore.First(user)
	expected := hex.EncodeToString(currentSum[:])
	if user.Password != expected {
		t.Errorf("password should not change; got %s", user.Password)
	}
}

func TestController_ProfileChangePassword_EmptyNewPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileController(t)
	currentSum := sha256.Sum256([]byte("oldpass"))
	user := &model.User{Username: "u", Password: hex.EncodeToString(currentSum[:])}
	ctrl.datastore.Create(user)

	body := `{"current_password":"oldpass","new_password":""}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/profile/password", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user", user)

	ctrl.ProfileChangePassword(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if errMsg, _ := resp["error"].(string); errMsg != "new password cannot be empty" {
		t.Errorf("expected error 'new password cannot be empty', got %q", errMsg)
	}
}

func TestController_ProfileChangePassword_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileController(t)
	xSum := sha256.Sum256([]byte("x"))
	user := &model.User{Username: "u", Password: hex.EncodeToString(xSum[:])}
	ctrl.datastore.Create(user)

	body := `{invalid json`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/profile/password", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user", user)

	ctrl.ProfileChangePassword(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if errMsg, _ := resp["error"].(string); errMsg != "invalid request" {
		t.Errorf("expected error 'invalid request', got %q", errMsg)
	}
}
