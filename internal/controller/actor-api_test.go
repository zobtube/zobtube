package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/provider"
)

// -----------------------------------------------------------------------------
// Setup helper
// -----------------------------------------------------------------------------

func setupActorController(t *testing.T) *Controller {
	cfg := &config.Config{
		DB: struct {
			Driver     string
			Connstring string
		}{
			Driver:     "sqlite",
			Connstring: ":memory:",
		},
	}
	db, err := model.New(cfg)
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	var logger zerolog.Logger
	if testing.Verbose() {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	} else {
		logger = zerolog.Nop()
	}
	shutdown := make(chan int, 1)
	controllerConfig := &config.Config{}
	ctrl := &Controller{
		logger:          &logger,
		shutdownChannel: shutdown,
		config:          controllerConfig,
		providers:       make(map[string]provider.Provider),
	}
	ctrl.DatabaseRegister(db)

	// create first user
	newUser := &model.User{
		Username: "admin",
		Admin:    true,
	}

	// save it
	err = ctrl.datastore.Save(&newUser).Error
	if err != nil {
		ctrl.logger.Error().Str("kind", "system").Err(err).Msg("unable to create initial user")
		t.Fatalf("failed to create initial user: %v", err)
	}

	// create configuration
	config := &model.Configuration{
		ID:                 1,
		UserAuthentication: false,
	}

	// save it
	err = ctrl.datastore.Save(&config).Error
	if err != nil {
		ctrl.logger.Error().Str("kind", "system").Err(err).Msg("unable to create initial user")
		t.Fatalf("failed to create configuration: %v", err)
	}
	ctrl.ConfigurationFromDBApply(config)

	return ctrl
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func TestController_ActorAPI_ActorAjaxNew_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupActorController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/actor/new", strings.NewReader(`name=Alice`))
	c.Request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	c.Request.ParseForm() //nolint:all
	c.Request.PostForm.Set("name", "Alice")

	ctrl.ActorAjaxNew(c)

	if w.Code != http.StatusOK {
		pageData, _ := io.ReadAll(w.Result().Body)
		t.Fatalf("expected 200, got %d, with error %s", w.Code, string(pageData))
	}

	actor := model.Actor{
		Name: "Alice",
	}

	result := ctrl.datastore.Find(&actor)
	if result.RowsAffected < 1 {
		t.Fatal("record not found")
	}
}

func TestController_ActorAPI_ActorAjaxProviderSearch_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupActorController(t)

	mockProv := &mockProvider{
		slug:              "mockprov",
		name:              "Mock Provider",
		searchActor:       true,
		scrapePicture:     false,
		actorSearchURL:    "https://mock.com/Alice",
		actorSearchResult: nil,
	}

	err := ctrl.ProviderRegister(mockProv)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// create actor item
	actor := model.Actor{
		Name: "Alice",
	}

	err = ctrl.datastore.Create(&actor).Error
	if err != nil {
		t.Fatalf("record not registered with error %s", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "provider_slug", Value: "mockprov"}}
	c.Request = httptest.NewRequest("GET", "/api/actor/"+actor.ID+"/provider/mockprov", nil)

	ctrl.ActorAjaxProviderSearch(c)

	if w.Code != http.StatusOK {
		pageData, _ := io.ReadAll(w.Result().Body)
		t.Fatalf("expected 200, got %d with error %s", w.Code, pageData)
	}

	var body map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	if body["link_url"] != "https://mock.com/Alice" {
		t.Errorf("expected url=https://mock.com/Alice, got %#v", body["link_url"])
	}
}
