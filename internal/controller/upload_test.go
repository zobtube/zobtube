package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

func setupUploadController(t *testing.T) *Controller {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&model.Video{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	logger := zerolog.Nop()
	shutdown := make(chan int, 1)
	ctrl := New(shutdown).(*Controller)
	ctrl.LoggerRegister(&logger)
	ctrl.DatabaseRegister(db)
	ctrl.ConfigurationRegister(&config.Config{})

	return ctrl
}

func TestController_UploadImport_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupUploadController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/upload/import", strings.NewReader(`{"path":"/media/triage/video.mp4","import_as":"v"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.UploadImport(c)

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
	if body["redirect"] == nil || body["redirect"] == "" {
		t.Error("expected redirect in response")
	}

	var vid model.Video
	if ctrl.datastore.First(&vid).RowsAffected < 1 {
		t.Error("video not created in DB")
	}
	if vid.Type != "v" {
		t.Errorf("expected type v, got %q", vid.Type)
	}
}

func TestController_UploadImport_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := setupUploadController(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/upload/import", strings.NewReader(`{invalid`))
	c.Request.Header.Set("Content-Type", "application/json")

	ctrl.UploadImport(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func setupAssignImageController(t *testing.T) (*Controller, string, string) {
	t.Helper()
	tmp := t.TempDir()
	libPath := filepath.Join(tmp, "lib")
	metaPath := filepath.Join(tmp, "meta")
	if err := os.MkdirAll(libPath, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(metaPath, 0o750); err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(
		&model.Video{}, &model.Actor{}, &model.Channel{}, &model.CategorySub{},
		&model.Library{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	defaultLibID, err := model.EnsureDefaultLibrary(db, libPath)
	if err != nil {
		t.Fatalf("failed to ensure default library: %v", err)
	}

	log := zerolog.Nop()
	shutdown := make(chan int, 1)
	ctrl := New(shutdown).(*Controller)
	ctrl.LoggerRegister(&log)
	ctrl.DatabaseRegister(db)
	cfg := &config.Config{DefaultLibraryID: defaultLibID}
	ctrl.ConfigurationRegister(cfg)
	resolver := storage.NewResolver(db)
	ctrl.StorageResolverRegister(resolver)
	ctrl.MetadataStorageRegister(storage.NewFilesystem(metaPath))

	r := &runner.Runner{}
	r.RegisterTask(&common.Task{
		Name:  "video/mini-thumb",
		Steps: []common.Step{{Name: "noop", NiceName: "No-op", Func: func(*common.Context, common.Parameters) (string, error) { return "", nil }}},
	})
	r.Start(cfg, db, resolver, storage.NewFilesystem(libPath))
	ctrl.RunnerRegister(r)

	return ctrl, libPath, metaPath
}

func writeTriageImage(t *testing.T, libPath, name string) {
	t.Helper()
	path := filepath.Join(libPath, "triage", name)
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte{0xff, 0xd8, 0xff, 0xd9}, 0o640); err != nil {
		t.Fatal(err)
	}
}

func TestController_UploadAssignImage_Actor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl, libPath, metaPath := setupAssignImageController(t)
	writeTriageImage(t, libPath, "photo.jpg")

	actor := &model.Actor{Name: "Test Actor"}
	if err := ctrl.datastore.Create(actor).Error; err != nil {
		t.Fatal(err)
	}

	body := `{"file":"photo.jpg","target_type":"actor","target_id":"` + actor.ID + `"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/upload/triage/assign-image", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	ctrl.UploadAssignImage(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp["redirect"] != "/actor/"+actor.ID+"/edit" {
		t.Errorf("redirect = %v", resp["redirect"])
	}

	var updated model.Actor
	if err := ctrl.datastore.First(&updated, "id = ?", actor.ID).Error; err != nil {
		t.Fatal(err)
	}
	if !updated.Thumbnail || !updated.Migrated {
		t.Errorf("expected thumbnail and migrated true, got thumbnail=%v migrated=%v", updated.Thumbnail, updated.Migrated)
	}

	thumbPath := filepath.Join(metaPath, "actors", actor.ID, "thumb.jpg")
	if _, err := os.Stat(thumbPath); err != nil {
		t.Fatalf("thumbnail not written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(libPath, "triage", "photo.jpg")); err != nil {
		t.Fatal("triage file should remain when delete_from_triage is false")
	}
}

func TestController_UploadAssignImage_DeleteFromTriage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl, libPath, _ := setupAssignImageController(t)
	writeTriageImage(t, libPath, "remove-me.jpg")

	channel := &model.Channel{Name: "Ch"}
	if err := ctrl.datastore.Create(channel).Error; err != nil {
		t.Fatal(err)
	}

	body := `{"file":"remove-me.jpg","target_type":"channel","target_id":"` + channel.ID + `","delete_from_triage":true}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/upload/triage/assign-image", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	ctrl.UploadAssignImage(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if _, err := os.Stat(filepath.Join(libPath, "triage", "remove-me.jpg")); !os.IsNotExist(err) {
		t.Fatal("expected triage file to be deleted")
	}
}

func TestController_UploadAssignImage_InvalidExtension(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl, libPath, _ := setupAssignImageController(t)
	path := filepath.Join(libPath, "triage", "doc.txt")
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("x"), 0o640); err != nil {
		t.Fatal(err)
	}

	actor := &model.Actor{Name: "A"}
	if err := ctrl.datastore.Create(actor).Error; err != nil {
		t.Fatal(err)
	}

	body := `{"file":"doc.txt","target_type":"actor","target_id":"` + actor.ID + `"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/upload/triage/assign-image", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	ctrl.UploadAssignImage(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
