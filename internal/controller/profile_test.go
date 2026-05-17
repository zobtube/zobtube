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

func setupProfileController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Video{}, &model.VideoView{}, &model.Actor{}); err != nil {
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

func TestController_ProfileView_OmitsOrphanVideoViews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileController(t)
	user := &model.User{Username: "u", Password: "x"}
	ctrl.datastore.Create(user)

	vid := &model.Video{Name: "exists", Filename: "exists.mp4", Type: "v"}
	ctrl.datastore.Create(vid)
	ctrl.datastore.Create(&model.VideoView{VideoID: vid.ID, UserID: user.ID, Count: 1})
	ctrl.datastore.Create(&model.VideoView{VideoID: "missing-video-id", UserID: user.ID, Count: 99})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/profile", nil)
	c.Set("user", user)

	ctrl.ProfileView(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	views, ok := resp["video_views"].([]any)
	if !ok {
		t.Fatalf("expected video_views array, got %T", resp["video_views"])
	}
	if len(views) != 1 {
		t.Fatalf("expected 1 video_view, got %d", len(views))
	}
	var remain int64
	ctrl.datastore.Model(&model.VideoView{}).Where("video_id = ?", "missing-video-id").Count(&remain)
	if remain != 0 {
		t.Errorf("expected orphan video_view removed, count=%d", remain)
	}
	stats := profileStatsFromResponse(t, resp)
	if stats["videos_unique"] != float64(1) {
		t.Errorf("stats videos_unique: got %v want 1", stats["videos_unique"])
	}
	if stats["videos_total"] != float64(1) {
		t.Errorf("stats videos_total: got %v want 1", stats["videos_total"])
	}
}

func profileStatsFromResponse(t *testing.T, resp map[string]any) map[string]any {
	t.Helper()
	stats, ok := resp["stats"].(map[string]any)
	if !ok {
		t.Fatalf("expected stats object, got %T", resp["stats"])
	}
	return stats
}

func TestController_ProfileView_Stats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileController(t)
	user := &model.User{Username: "u", Password: "x"}
	ctrl.datastore.Create(user)

	d1 := 10 * time.Minute
	d2 := 5 * time.Minute
	v1 := &model.Video{Name: "v1", Filename: "v1.mp4", Type: "v", Duration: d1}
	v2 := &model.Video{Name: "v2", Filename: "v2.mp4", Type: "v", Duration: d2}
	actor := &model.Actor{Name: "Star"}
	ctrl.datastore.Create(v1)
	ctrl.datastore.Create(v2)
	ctrl.datastore.Create(actor)
	if err := ctrl.datastore.Model(v1).Association("Actors").Append(actor); err != nil {
		t.Fatalf("link actor v1: %v", err)
	}
	if err := ctrl.datastore.Model(v2).Association("Actors").Append(actor); err != nil {
		t.Fatalf("link actor v2: %v", err)
	}
	ctrl.datastore.Create(&model.VideoView{VideoID: v1.ID, UserID: user.ID, Count: 2})
	ctrl.datastore.Create(&model.VideoView{VideoID: v2.ID, UserID: user.ID, Count: 3})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/profile", nil)
	c.Set("user", user)
	ctrl.ProfileView(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	stats := profileStatsFromResponse(t, resp)
	if stats["videos_unique"] != float64(2) {
		t.Errorf("videos_unique: got %v want 2", stats["videos_unique"])
	}
	if stats["videos_total"] != float64(5) {
		t.Errorf("videos_total: got %v want 5", stats["videos_total"])
	}
	if stats["actors_unique"] != float64(1) {
		t.Errorf("actors_unique: got %v want 1", stats["actors_unique"])
	}
	if stats["actors_total"] != float64(5) {
		t.Errorf("actors_total: got %v want 5", stats["actors_total"])
	}
	wantTime := int64(2)*int64(d1) + int64(3)*int64(d2)
	if got, ok := stats["total_view_time_ns"].(float64); !ok || int64(got) != wantTime {
		t.Errorf("total_view_time_ns: got %v want %d", stats["total_view_time_ns"], wantTime)
	}
}

func TestController_ProfileView_MigratesViewsFromDeletedVideo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupProfileController(t)
	user := &model.User{Username: "u", Password: "x"}
	ctrl.datastore.Create(user)
	libID := "lib-1"

	deleted := &model.Video{Name: "old", Filename: "shared.mp4", Type: "v", LibraryID: &libID}
	ctrl.datastore.Create(deleted)
	ctrl.datastore.Delete(deleted)

	replacement := &model.Video{Name: "movie 1", Filename: "shared.mp4", Type: "m", LibraryID: &libID}
	ctrl.datastore.Create(replacement)
	ctrl.datastore.Create(&model.VideoView{VideoID: deleted.ID, UserID: user.ID, Count: 4})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/profile", nil)
	c.Set("user", user)
	ctrl.ProfileView(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var migrated model.VideoView
	if err := ctrl.datastore.First(&migrated, "video_id = ? AND user_id = ?", replacement.ID, user.ID).Error; err != nil {
		t.Fatalf("expected migrated view on replacement: %v", err)
	}
	if migrated.Count != 4 {
		t.Errorf("expected count 4 on replacement, got %d", migrated.Count)
	}
	var oldRemain int64
	ctrl.datastore.Model(&model.VideoView{}).Where("video_id = ?", deleted.ID).Count(&oldRemain)
	if oldRemain != 0 {
		t.Errorf("expected old video_view removed, count=%d", oldRemain)
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	stats := profileStatsFromResponse(t, resp)
	if stats["videos_unique"] != float64(1) {
		t.Errorf("stats videos_unique: got %v want 1", stats["videos_unique"])
	}
	if stats["videos_total"] != float64(4) {
		t.Errorf("stats videos_total: got %v want 4", stats["videos_total"])
	}
}
