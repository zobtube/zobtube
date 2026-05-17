package controller

import (
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

func setupPlaylistController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Video{}, &model.Playlist{}, &model.PlaylistVideo{}); err != nil {
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

func TestController_PlaylistList_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupPlaylistController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/playlists", nil)
	c.Set("user", user)

	ctrl.PlaylistList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	playlists, _ := resp["playlists"].([]any)
	if len(playlists) != 0 {
		t.Errorf("expected 0 playlists, got %d", len(playlists))
	}
}

func TestController_PlaylistCreate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupPlaylistController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/playlists", strings.NewReader(`{"name":"Favorites"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user", user)

	ctrl.PlaylistCreate(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["name"] != "Favorites" {
		t.Errorf("expected name Favorites, got %v", resp["name"])
	}
}

func TestController_PlaylistCreate_EmptyName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupPlaylistController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/playlists", strings.NewReader(`{"name":""}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user", user)

	ctrl.PlaylistCreate(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestController_PlaylistView_NotFoundWrongUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupPlaylistController(t)
	user1 := &model.User{Username: "u1"}
	user2 := &model.User{Username: "u2"}
	ctrl.datastore.Create(user1)
	ctrl.datastore.Create(user2)
	now := time.Now()
	pl := &model.Playlist{UserID: user1.ID, Name: "Mine", CreatedAt: now, UpdatedAt: now}
	ctrl.datastore.Create(pl)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/playlists/"+pl.ID, nil)
	c.Params = gin.Params{{Key: "id", Value: pl.ID}}
	c.Set("user", user2)

	ctrl.PlaylistView(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestController_PlaylistVideoAdd_Remove(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupPlaylistController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)
	now := time.Now()
	pl := &model.Playlist{UserID: user.ID, Name: "List", CreatedAt: now, UpdatedAt: now}
	ctrl.datastore.Create(pl)
	video := &model.Video{Name: "V", Type: "v", Status: model.VideoStatusReady, CreatedAt: now, UpdatedAt: now}
	ctrl.datastore.Create(video)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/playlists/"+pl.ID+"/videos", strings.NewReader(`{"video_id":"`+video.ID+`"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: pl.ID}}
	c.Set("user", user)

	ctrl.PlaylistVideoAdd(c)
	if w.Code != http.StatusOK {
		t.Fatalf("add: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest("GET", "/api/playlists/"+pl.ID, nil)
	c2.Params = gin.Params{{Key: "id", Value: pl.ID}}
	c2.Set("user", user)
	ctrl.PlaylistView(c2)
	if w2.Code != http.StatusOK {
		t.Fatalf("view: expected 200, got %d", w2.Code)
	}
	var viewResp map[string]any
	json.Unmarshal(w2.Body.Bytes(), &viewResp)
	videos, _ := viewResp["videos"].([]any)
	if len(videos) != 1 {
		t.Fatalf("expected 1 video, got %d", len(videos))
	}

	w3 := httptest.NewRecorder()
	c3, _ := gin.CreateTestContext(w3)
	c3.Request = httptest.NewRequest("DELETE", "/api/playlists/"+pl.ID+"/videos/"+video.ID, nil)
	c3.Params = gin.Params{{Key: "id", Value: pl.ID}, {Key: "video_id", Value: video.ID}}
	c3.Set("user", user)
	ctrl.PlaylistVideoRemove(c3)
	if w3.Code != http.StatusNoContent {
		t.Fatalf("remove: expected 204, got %d: %s", w3.Code, w3.Body.String())
	}
}

func TestController_PlaylistList_ContainsVideo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupPlaylistController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)
	now := time.Now()
	pl := &model.Playlist{UserID: user.ID, Name: "List", CreatedAt: now, UpdatedAt: now}
	ctrl.datastore.Create(pl)
	video := &model.Video{Name: "V", Type: "v", Status: model.VideoStatusReady, CreatedAt: now, UpdatedAt: now}
	ctrl.datastore.Create(video)
	ctrl.datastore.Create(&model.PlaylistVideo{PlaylistID: pl.ID, VideoID: video.ID, Position: 0, AddedAt: now})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/playlists?video_id="+video.ID, nil)
	c.Set("user", user)

	ctrl.PlaylistList(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	playlists, _ := resp["playlists"].([]any)
	if len(playlists) != 1 {
		t.Fatalf("expected 1 playlist, got %d", len(playlists))
	}
	item, _ := playlists[0].(map[string]any)
	if item["contains"] != true {
		t.Errorf("expected contains true, got %v", item["contains"])
	}
}

func TestController_PlaylistPlaybackContext(t *testing.T) {
	ctrl := setupPlaylistController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)
	now := time.Now()
	pl := &model.Playlist{UserID: user.ID, Name: "Queue", CreatedAt: now, UpdatedAt: now}
	ctrl.datastore.Create(pl)
	v1 := &model.Video{Name: "V1", Type: "v", Status: model.VideoStatusReady, CreatedAt: now, UpdatedAt: now}
	v2 := &model.Video{Name: "V2", Type: "v", Status: model.VideoStatusReady, CreatedAt: now, UpdatedAt: now}
	v3 := &model.Video{Name: "V3", Type: "v", Status: model.VideoStatusReady, CreatedAt: now, UpdatedAt: now}
	ctrl.datastore.Create(v1)
	ctrl.datastore.Create(v2)
	ctrl.datastore.Create(v3)
	ctrl.datastore.Create(&model.PlaylistVideo{PlaylistID: pl.ID, VideoID: v1.ID, Position: 0, AddedAt: now})
	ctrl.datastore.Create(&model.PlaylistVideo{PlaylistID: pl.ID, VideoID: v2.ID, Position: 1, AddedAt: now})
	ctrl.datastore.Create(&model.PlaylistVideo{PlaylistID: pl.ID, VideoID: v3.ID, Position: 2, AddedAt: now})

	ctx := ctrl.playlistPlaybackContext(user.ID, pl.ID, v2.ID)
	if ctx == nil {
		t.Fatal("expected playback context")
	}
	if ctx["playlist_index"] != 1 {
		t.Errorf("expected index 1, got %v", ctx["playlist_index"])
	}
	ids, _ := ctx["playlist_video_ids"].([]string)
	if len(ids) != 3 || ids[0] != v1.ID || ids[2] != v3.ID {
		t.Errorf("unexpected ids: %v", ids)
	}
	upNext, _ := ctx["playlist_up_next"].([]model.Video)
	if len(upNext) != 1 || upNext[0].ID != v3.ID {
		t.Errorf("expected up_next [v3], got %v", upNext)
	}
	items, _ := ctx["playlist_items"].([]gin.H)
	if len(items) != 3 {
		t.Fatalf("expected 3 playlist_items, got %d", len(items))
	}
	if items[1]["id"] != v2.ID || items[1]["type"] != "v" {
		t.Errorf("unexpected playlist_items[1]: %v", items[1])
	}
	allVideos, _ := ctx["playlist_videos"].([]model.Video)
	if len(allVideos) != 3 {
		t.Fatalf("expected 3 playlist_videos, got %d", len(allVideos))
	}
	if allVideos[0].ID != v1.ID || allVideos[1].ID != v2.ID || allVideos[2].ID != v3.ID {
		t.Errorf("unexpected playlist_videos order: %v, %v, %v", allVideos[0].ID, allVideos[1].ID, allVideos[2].ID)
	}
	if ctrl.playlistPlaybackContext(user.ID, pl.ID, "00000000-0000-0000-0000-000000000000") != nil {
		t.Error("expected nil for video not in playlist")
	}
	user2 := &model.User{Username: "u2"}
	ctrl.datastore.Create(user2)
	if ctrl.playlistPlaybackContext(user2.ID, pl.ID, v1.ID) != nil {
		t.Error("expected nil for wrong user")
	}
}

func TestController_PlaylistDelete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupPlaylistController(t)
	user := &model.User{Username: "u"}
	ctrl.datastore.Create(user)
	now := time.Now()
	pl := &model.Playlist{UserID: user.ID, Name: "List", CreatedAt: now, UpdatedAt: now}
	ctrl.datastore.Create(pl)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/api/playlists/"+pl.ID, nil)
	c.Params = gin.Params{{Key: "id", Value: pl.ID}}
	c.Set("user", user)

	ctrl.PlaylistDelete(c)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
	var count int64
	ctrl.datastore.Model(&model.Playlist{}).Where("id = ?", pl.ID).Count(&count)
	if count != 0 {
		t.Error("playlist should be deleted")
	}
}
