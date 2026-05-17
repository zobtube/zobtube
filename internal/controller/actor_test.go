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

func TestController_Actor_ActorNew_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupActorController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/actor/new", strings.NewReader(`name=Alice`))
	c.Request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	c.Request.ParseForm() //nolint:all
	c.Request.PostForm.Set("name", "Alice")

	ctrl.ActorNew(c)

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

func TestController_ActorList_SortedByNameCaseInsensitive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupActorController(t)

	for _, name := range []string{"Zara", "alice", "Bob"} {
		if err := ctrl.datastore.Create(&model.Actor{Name: name}).Error; err != nil {
			t.Fatalf("create actor %q: %v", name, err)
		}
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/actor", nil)
	ctrl.ActorList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Items []model.Actor `json:"items"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Items) != 3 {
		t.Fatalf("expected 3 actors, got %d", len(resp.Items))
	}

	got := []string{resp.Items[0].Name, resp.Items[1].Name, resp.Items[2].Name}
	want := []string{"alice", "Bob", "Zara"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order[%d]: got %q, want %q (full order: %v)", i, got[i], want[i], got)
		}
	}
}

func TestController_Actor_ActorProviderSearch_Success(t *testing.T) {
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

	ctrl.ActorProviderSearch(c)

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

func TestController_Actor_ActorMerge_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupActorController(t)

	source := &model.Actor{Name: "Source Actor", Sex: "f"}
	target := &model.Actor{Name: "Target Actor", Sex: "f"}
	if err := ctrl.datastore.Create(source).Error; err != nil {
		t.Fatalf("create source: %v", err)
	}
	if err := ctrl.datastore.Create(target).Error; err != nil {
		t.Fatalf("create target: %v", err)
	}

	video := &model.Video{Name: "Test Video", Filename: "test.mp4", Type: "v"}
	if err := ctrl.datastore.Create(video).Error; err != nil {
		t.Fatalf("create video: %v", err)
	}
	if err := ctrl.datastore.Model(video).Association("Actors").Append(source); err != nil {
		t.Fatalf("link video to source: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: source.ID}}
	c.Request = httptest.NewRequest("POST", "/api/actor/"+source.ID+"/merge", strings.NewReader(`{"target_id":"`+target.ID+`"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.ActorMerge(c)

	if w.Code != http.StatusOK {
		body, _ := io.ReadAll(w.Result().Body)
		t.Fatalf("expected 200, got %d: %s", w.Code, string(body))
	}

	var videoCheck model.Video
	if err := ctrl.datastore.Preload("Actors").First(&videoCheck, "id = ?", video.ID).Error; err != nil {
		t.Fatalf("load video: %v", err)
	}
	if len(videoCheck.Actors) != 1 || videoCheck.Actors[0].ID != target.ID {
		t.Errorf("video should have only target actor; got %d actors: %v", len(videoCheck.Actors), videoCheck.Actors)
	}

	var sourceCheck model.Actor
	if res := ctrl.datastore.First(&sourceCheck, "id = ?", source.ID); res.RowsAffected > 0 {
		t.Error("source actor should be deleted (not found by normal query)")
	}
}

func TestController_Actor_ActorMerge_SameID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupActorController(t)

	actor := &model.Actor{Name: "Only", Sex: "f"}
	if err := ctrl.datastore.Create(actor).Error; err != nil {
		t.Fatalf("create actor: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: actor.ID}}
	c.Request = httptest.NewRequest("POST", "/api/actor/"+actor.ID+"/merge", strings.NewReader(`{"target_id":"`+actor.ID+`"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.ActorMerge(c)

	if w.Code != http.StatusBadRequest {
		body, _ := io.ReadAll(w.Result().Body)
		t.Fatalf("expected 400, got %d: %s", w.Code, string(body))
	}
}

func TestController_ActorPhotosets_UnionExcludesDeleting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupActorController(t)

	actor := &model.Actor{Name: "Tagged Actor"}
	if err := ctrl.datastore.Create(actor).Error; err != nil {
		t.Fatalf("create actor: %v", err)
	}

	psAlbum := &model.Photoset{Name: "Album tag", Status: model.PhotosetStatusReady}
	psPhoto := &model.Photoset{Name: "Photo tag", Status: model.PhotosetStatusReady}
	psDeleting := &model.Photoset{Name: "Deleting", Status: model.PhotosetStatusDeleting}
	for _, ps := range []*model.Photoset{psAlbum, psPhoto, psDeleting} {
		if err := ctrl.datastore.Create(ps).Error; err != nil {
			t.Fatalf("create photoset: %v", err)
		}
	}
	if err := ctrl.datastore.Model(psAlbum).Association("Actors").Append(actor); err != nil {
		t.Fatalf("album actor link: %v", err)
	}
	if err := ctrl.datastore.Model(psDeleting).Association("Actors").Append(actor); err != nil {
		t.Fatalf("deleting album actor link: %v", err)
	}

	photo := &model.Photo{PhotosetID: psPhoto.ID, Filename: "one.jpg"}
	if err := ctrl.datastore.Create(photo).Error; err != nil {
		t.Fatalf("create photo: %v", err)
	}
	if err := ctrl.datastore.Model(photo).Association("Actors").Append(actor); err != nil {
		t.Fatalf("photo actor link: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: actor.ID}}
	c.Request = httptest.NewRequest("GET", "/api/actor/"+actor.ID+"/photosets", nil)
	ctrl.ActorPhotosets(c)

	if w.Code != http.StatusOK {
		body, _ := io.ReadAll(w.Result().Body)
		t.Fatalf("expected 200, got %d: %s", w.Code, string(body))
	}

	var resp struct {
		Items []model.Photoset `json:"items"`
		Total int              `json:"total"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Total != 2 || len(resp.Items) != 2 {
		t.Fatalf("expected 2 photosets, got total=%d len=%d", resp.Total, len(resp.Items))
	}
	ids := map[string]bool{resp.Items[0].ID: true, resp.Items[1].ID: true}
	if !ids[psAlbum.ID] || !ids[psPhoto.ID] {
		t.Fatalf("unexpected photoset ids: %v", ids)
	}
}
