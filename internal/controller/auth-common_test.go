package controller

import (
	"net/http/httptest"
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

// --- setup helpers ---

func setupAuthTestController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.UserSession{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	_logger := zerolog.Nop()
	cfg := &config.Config{Authentication: true}
	shutdown := make(chan int, 1)

	ctrl := New(shutdown).(*Controller)
	ctrl.LoggerRegister(&_logger)
	ctrl.DatabaseRegister(db)
	ctrl.ConfigurationRegister(cfg)

	return ctrl
}

// --- tests ---

func TestController_CreateSession_SetsCookieAndDBRecord(t *testing.T) {
	ctrl := setupAuthTestController(t)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ctrl.createSession(c)

	// Check cookie set
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected cookie to be set, got none")
	}
	found := false
	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			found = true
			if cookie.MaxAge <= 0 {
				t.Errorf("expected positive MaxAge, got %d", cookie.MaxAge)
			}
		}
	}
	if !found {
		t.Errorf("expected cookie %q to be set", cookieName)
	}

	// Check DB entry
	var count int64
	if err := ctrl.datastore.Model(&model.UserSession{}).Count(&count).Error; err != nil {
		t.Fatalf("db query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 session in DB, got %d", count)
	}
}

func TestController_GetSession(t *testing.T) {
	ctrl := setupAuthTestController(t)

	session := &model.UserSession{ValidUntil: time.Now().Add(10 * time.Minute)}
	if err := ctrl.datastore.Create(session).Error; err != nil {
		t.Fatalf("failed to insert session: %v", err)
	}

	var s model.UserSession
	ctrl.GetSession(&s)
	if s.ID == "" {
		t.Error("expected GetSession to load a record, got empty ID")
	}
}

func TestController_GetUser(t *testing.T) {
	ctrl := setupAuthTestController(t)

	user := &model.User{Username: "tester"}
	if err := ctrl.datastore.Create(user).Error; err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	var u model.User
	ctrl.GetUser(&u)
	if u.Username == "" {
		t.Error("expected GetUser to populate a user, got empty name")
	}
}

func TestController_GetFirstUser_ReturnsOldest(t *testing.T) {
	ctrl := setupAuthTestController(t)

	u1 := &model.User{Username: "first"}
	u2 := &model.User{Username: "second"}
	if err := ctrl.datastore.Create(u1).Error; err != nil {
		t.Fatalf("insert error: %v", err)
	}
	time.Sleep(time.Millisecond) // ensure ordering
	if err := ctrl.datastore.Create(u2).Error; err != nil {
		t.Fatalf("insert error: %v", err)
	}

	var u model.User
	ctrl.GetFirstUser(&u)
	if u.Username != "first" {
		t.Errorf("expected oldest user 'first', got %q", u.Username)
	}
}

func TestController_AuthenticationEnabled(t *testing.T) {
	ctrl := setupAuthTestController(t)

	// Should reflect config.Authentication
	if !ctrl.AuthenticationEnabled() {
		t.Error("expected AuthenticationEnabled() true, got false")
	}

	// Toggle and verify
	ctrl.config.Authentication = false
	if ctrl.AuthenticationEnabled() {
		t.Error("expected AuthenticationEnabled() false after toggle")
	}
}
