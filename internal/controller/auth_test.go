package controller

import (
	"crypto/sha256"
	"encoding/hex"
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

func setupAuthController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.UserSession{}); err != nil {
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

func TestController_AuthLogin_NoCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAuthController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", strings.NewReader(`username=u&password=p`))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request.ParseForm() //nolint:errcheck

	ctrl.AuthLogin(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestController_AuthLogin_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAuthController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", strings.NewReader(`username=nonexistent&password=x`))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request.ParseForm() //nolint:errcheck

	// Create a valid session and attach its cookie
	session := &model.UserSession{ValidUntil: time.Now().Add(time.Hour)}
	ctrl.datastore.Create(session)
	c.Request.AddCookie(&http.Cookie{Name: cookieName, Value: session.ID})

	ctrl.AuthLogin(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
	}
}

func TestController_AuthLogin_WrongPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAuthController(t)
	user := &model.User{Username: "u", Password: "secret"}
	ctrl.datastore.Create(user)
	session := &model.UserSession{ValidUntil: time.Now().Add(time.Hour)}
	ctrl.datastore.Create(session)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", strings.NewReader(`username=u&password=wrong`))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request.ParseForm() //nolint:errcheck
	c.Request.AddCookie(&http.Cookie{Name: cookieName, Value: session.ID})

	ctrl.AuthLogin(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestController_AuthLogin_ExpiredSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAuthController(t)
	user := &model.User{Username: "u", Password: "p"}
	ctrl.datastore.Create(user)
	session := &model.UserSession{ValidUntil: time.Now().Add(-time.Hour)}
	ctrl.datastore.Create(session)

	challengeHex := sha256.Sum256([]byte(session.ID + user.Password))
	challenge := hex.EncodeToString(challengeHex[:])

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/auth/login", strings.NewReader("username=u&password="+challenge))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request.ParseForm() //nolint:errcheck
	c.Request.AddCookie(&http.Cookie{Name: cookieName, Value: session.ID})

	ctrl.AuthLogin(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired session, got %d", w.Code)
	}
}

func TestController_AuthLogout_WithCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAuthController(t)
	session := &model.UserSession{ValidUntil: time.Now().Add(time.Hour)}
	ctrl.datastore.Create(session)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/auth/logout", nil)
	c.Request.AddCookie(&http.Cookie{Name: cookieName, Value: session.ID})

	ctrl.AuthLogout(c)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	// Session should be deleted
	var count int64
	ctrl.datastore.Model(&model.UserSession{}).Where("id = ?", session.ID).Count(&count)
	if count > 0 {
		t.Error("expected session to be deleted")
	}
}

func TestController_AuthLogout_WithoutCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAuthController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/auth/logout", nil)

	ctrl.AuthLogout(c)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestController_AuthLogoutRedirect_ClearsAndRedirects(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAuthController(t)
	session := &model.UserSession{ValidUntil: time.Now().Add(time.Hour)}
	ctrl.datastore.Create(session)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/auth/logout", nil)
	c.Request.AddCookie(&http.Cookie{Name: cookieName, Value: session.ID})

	ctrl.AuthLogoutRedirect(c)

	if w.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "/" {
		t.Errorf("expected Location /, got %s", loc)
	}
}
