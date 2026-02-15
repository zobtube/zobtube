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

func setupAppController(t *testing.T) *Controller {
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

	return ctrl
}

func TestController_Bootstrap_AuthDisabled_ReturnsFirstUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAppController(t)
	ctrl.ConfigurationRegister(&config.Config{Authentication: false})
	user := &model.User{Username: "testuser"}
	ctrl.datastore.Create(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/bootstrap", nil)

	ctrl.Bootstrap(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["auth_enabled"] != false {
		t.Errorf("expected auth_enabled false, got %v", body["auth_enabled"])
	}
	if body["user"] == nil {
		t.Error("expected user in response")
	}
}

func TestController_Bootstrap_AuthEnabled_NoCookie_CreatesSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAppController(t)
	ctrl.ConfigurationRegister(&config.Config{Authentication: true})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/bootstrap", nil)

	ctrl.Bootstrap(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	// Session should be created (cookie set)
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected session cookie to be set")
	}
}

func TestController_NoRouteOrSPA_ApiPath_Returns404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupAppController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/nonexistent", nil)

	ctrl.NoRouteOrSPA(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for API path, got %d", w.Code)
	}
}
