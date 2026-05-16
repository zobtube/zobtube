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

func setupUploadScanController(t *testing.T) (*Controller, string, string) {
	t.Helper()
	root := t.TempDir()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(
		&model.Video{}, &model.Actor{}, &model.Channel{}, &model.Category{},
		&model.CategorySub{}, &model.Task{}, &model.Library{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	libID, err := model.EnsureDefaultLibrary(db, root)
	if err != nil {
		t.Fatalf("default library: %v", err)
	}
	var lib model.Library
	if err := db.First(&lib, "id = ?", libID).Error; err != nil {
		t.Fatalf("load library: %v", err)
	}
	lib.Config = model.LibraryConfig{Filesystem: &model.LibraryConfigFilesystem{Path: root}}
	if err := db.Save(&lib).Error; err != nil {
		t.Fatalf("save library config: %v", err)
	}

	stageScanFixtures(t, root)

	log := zerolog.Nop()
	shutdown := make(chan int, 1)
	ctrl := New(shutdown).(*Controller)
	ctrl.LoggerRegister(&log)
	ctrl.DatabaseRegister(db)
	cfg := &config.Config{DefaultLibraryID: libID}
	ctrl.ConfigurationRegister(cfg)
	resolver := storage.NewResolver(db)
	ctrl.StorageResolverRegister(resolver)

	r := &runner.Runner{}
	r.RegisterTask(&common.Task{
		Name: "video/create",
		Steps: []common.Step{{
			Name: "noop", NiceName: "No-op",
			Func: func(*common.Context, common.Parameters) (string, error) { return "", nil },
		}},
	})
	r.Start(cfg, db, resolver, storage.NewFilesystem(root))
	ctrl.RunnerRegister(r)

	return ctrl, libID, root
}

func stageScanFixtures(t *testing.T, root string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "triage", "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeSizedFile(t, filepath.Join(root, "triage", "a.mp4"), 512*1024)       // 0 MB -> clip @ 1MB threshold
	writeSizedFile(t, filepath.Join(root, "triage", "sub", "b.mp4"), 3*1024*1024/2) // 1 MB -> video
	writeSizedFile(t, filepath.Join(root, "triage", "sub", "c.mp4"), 3*1024*1024)   // 3 MB -> movie
	if err := os.WriteFile(filepath.Join(root, "triage", "sub", "notes.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeSizedFile(t *testing.T, path string, size int64) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(size); err != nil {
		f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

func postScan(t *testing.T, ctrl *Controller, libID string, body string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/upload/triage/scan", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	ctrl.UploadTriageScan(c)
	return w
}

func scanBody(libID string, recursive bool, enC, enV, enM bool, clipMB, videoMB int64) string {
	b, _ := json.Marshal(map[string]any{
		"path":      "",
		"recursive": recursive,
		"library_id": libID,
		"enabled": map[string]bool{"c": enC, "v": enV, "m": enM},
		"thresholds": map[string]int64{
			"clip_video_mb":  clipMB,
			"video_movie_mb": videoMB,
		},
	})
	return string(b)
}

func TestUploadTriageScan_NonRecursive(t *testing.T) {
	ctrl, libID, _ := setupUploadScanController(t)
	w := postScan(t, ctrl, libID, scanBody(libID, false, true, true, true, 1, 2))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if int(resp["imported"].(float64)) != 1 {
		t.Fatalf("imported = %v, want 1", resp["imported"])
	}
	var vids []model.Video
	if err := ctrl.datastore.Find(&vids).Error; err != nil {
		t.Fatal(err)
	}
	if len(vids) != 1 || vids[0].Filename != "a.mp4" || vids[0].Type != "c" {
		t.Fatalf("videos = %+v", vids)
	}
}

func TestUploadTriageScan_RecursiveClassification(t *testing.T) {
	ctrl, libID, _ := setupUploadScanController(t)
	w := postScan(t, ctrl, libID, scanBody(libID, true, true, true, true, 1, 2))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["imported"].(float64)) != 3 {
		t.Fatalf("imported = %v, want 3", resp["imported"])
	}
	types := map[string]string{}
	var vids []model.Video
	if err := ctrl.datastore.Find(&vids).Error; err != nil {
		t.Fatal(err)
	}
	for _, v := range vids {
		types[v.Filename] = v.Type
	}
	if types["a.mp4"] != "c" || types["sub/b.mp4"] != "v" || types["sub/c.mp4"] != "m" {
		t.Fatalf("types = %v", types)
	}
	if int(resp["skipped_non_video"].(float64)) != 1 {
		t.Fatalf("skipped_non_video = %v, want 1", resp["skipped_non_video"])
	}
}

func TestUploadTriageScan_DisabledType(t *testing.T) {
	ctrl, libID, _ := setupUploadScanController(t)
	w := postScan(t, ctrl, libID, scanBody(libID, true, true, false, true, 1, 2))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if int(resp["skipped_disabled"].(float64)) < 1 {
		t.Fatalf("skipped_disabled = %v, want >= 1", resp["skipped_disabled"])
	}
	var v model.Video
	if ctrl.datastore.Where("filename = ?", "sub/b.mp4").First(&v).RowsAffected > 0 {
		t.Fatal("sub/b.mp4 should not be imported when video type disabled")
	}
}

func TestUploadTriageScan_SkipsExisting(t *testing.T) {
	ctrl, libID, _ := setupUploadScanController(t)
	body := scanBody(libID, true, true, true, true, 1, 2)
	w1 := postScan(t, ctrl, libID, body)
	if w1.Code != http.StatusOK {
		t.Fatalf("first scan: %d %s", w1.Code, w1.Body.String())
	}
	w2 := postScan(t, ctrl, libID, body)
	if w2.Code != http.StatusOK {
		t.Fatalf("second scan: %d %s", w2.Code, w2.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(w2.Body.Bytes(), &resp)
	if int(resp["imported"].(float64)) != 0 {
		t.Fatalf("imported = %v, want 0", resp["imported"])
	}
	if int(resp["skipped_existing"].(float64)) != 3 {
		t.Fatalf("skipped_existing = %v, want 3", resp["skipped_existing"])
	}
}

func TestUploadTriageScan_Validation(t *testing.T) {
	ctrl, libID, _ := setupUploadScanController(t)

	cases := []struct {
		name string
		body string
	}{
		{"no types", `{"path":"","enabled":{"c":false,"v":false,"m":false},"thresholds":{"clip_video_mb":1,"video_movie_mb":2}}`},
		{"bad thresholds", `{"path":"","enabled":{"c":true,"v":true,"m":true},"thresholds":{"clip_video_mb":5,"video_movie_mb":2}}`},
		{"missing actor", `{"path":"","enabled":{"c":true,"v":true,"m":true},"thresholds":{"clip_video_mb":1,"video_movie_mb":2},"actors":["missing"]}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := postScan(t, ctrl, libID, tc.body)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
			}
		})
	}
}
