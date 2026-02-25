package model

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupVideoDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&Video{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestVideo_BeforeCreate_AssignsUUIDWhenIDEmpty(t *testing.T) {
	db := setupVideoDB(t)
	v := &Video{Name: "test", Type: "v"}
	if err := db.Create(v).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if v.ID == "" {
		t.Error("expected ID to be set by BeforeCreate")
	}
	if _, err := uuid.Parse(v.ID); err != nil {
		t.Errorf("expected ID to be valid UUID, got %q: %v", v.ID, err)
	}
}

func TestVideo_BeforeCreate_AssignsUUIDWhenZeroUUID(t *testing.T) {
	db := setupVideoDB(t)
	v := &Video{ID: "00000000-0000-0000-0000-000000000000", Name: "test", Type: "v"}
	if err := db.Create(v).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if v.ID == "00000000-0000-0000-0000-000000000000" {
		t.Error("expected zero UUID to be replaced")
	}
	if _, err := uuid.Parse(v.ID); err != nil {
		t.Errorf("expected ID to be valid UUID, got %q: %v", v.ID, err)
	}
}

func TestVideo_BeforeCreate_KeepsIDWhenSet(t *testing.T) {
	db := setupVideoDB(t)
	wantID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	v := &Video{ID: wantID, Name: "test", Type: "v"}
	if err := db.Create(v).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if v.ID != wantID {
		t.Errorf("expected ID to remain %q, got %q", wantID, v.ID)
	}
}

func TestVideo_TypeAsString(t *testing.T) {
	tests := []struct {
		typeChar string
		want     string
	}{
		{"c", "clip"},
		{"v", "video"},
		{"m", "movie"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.typeChar, func(t *testing.T) {
			v := &Video{Type: tt.typeChar}
			got := v.TypeAsString()
			if got != tt.want {
				t.Errorf("TypeAsString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestVideo_URLView(t *testing.T) {
	id := "deadbeef-1234-5678-abcd-000000000000"
	t.Run("clip", func(t *testing.T) {
		v := &Video{ID: id, Type: "c"}
		got := v.URLView()
		want := "/clip/" + id
		if got != want {
			t.Errorf("URLView() = %q, want %q", got, want)
		}
	})
	t.Run("video", func(t *testing.T) {
		v := &Video{ID: id, Type: "v"}
		got := v.URLView()
		want := "/video/" + id
		if got != want {
			t.Errorf("URLView() = %q, want %q", got, want)
		}
	})
}

func TestVideo_URLThumb_URLThumbXS_URLStream_URLAdmEdit(t *testing.T) {
	id := "deadbeef-1234-5678-abcd-000000000000"
	v := &Video{ID: id}
	if got := v.URLThumb(); got != "/api/video/"+id+"/thumb" {
		t.Errorf("URLThumb() = %q", got)
	}
	if got := v.URLThumbXS(); got != "/api/video/"+id+"/thumb_xs" {
		t.Errorf("URLThumbXS() = %q", got)
	}
	if got := v.URLStream(); got != "/api/video/"+id+"/stream" {
		t.Errorf("URLStream() = %q", got)
	}
	if got := v.URLAdmEdit(); got != "/video/"+id+"/edit" {
		t.Errorf("URLAdmEdit() = %q", got)
	}
}

func TestVideo_NiceDuration(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{"zero", 0, "00:00"},
		{"90s", 90 * time.Second, "01:30"},
		{"1h1m5s", time.Hour + time.Minute + 5*time.Second, "01:01:05"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Video{Duration: tt.d}
			got := v.NiceDuration()
			if got != tt.want {
				t.Errorf("NiceDuration() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestVideo_NiceDurationShort(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{"0s", 0, " 0 sec"},
		{"90s", 90 * time.Second, " 1 min"},
		{"1h1m", time.Hour + time.Minute, " 1h01"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Video{Duration: tt.d}
			got := v.NiceDurationShort()
			if got != tt.want {
				t.Errorf("NiceDurationShort() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestVideo_String(t *testing.T) {
	v := &Video{Name: "My Video"}
	if got := v.String(); got != "My Video" {
		t.Errorf("String() = %q, want %q", got, "My Video")
	}
}

func TestVideo_HasDuration(t *testing.T) {
	v := &Video{}
	if !v.HasDuration() {
		t.Error("HasDuration() = false, want true")
	}
}

func TestVideo_FolderRelativePath_RelativePath_ThumbPaths(t *testing.T) {
	id := "deadbeef-1234-5678-abcd-000000000000"
	t.Run("clip", func(t *testing.T) {
		v := &Video{ID: id, Type: "c"}
		if got := v.FolderRelativePath(); got != filepath.Join("/clips", id) {
			t.Errorf("FolderRelativePath() = %q", got)
		}
		if got := v.RelativePath(); got != filepath.Join("/clips", id, "video.mp4") {
			t.Errorf("RelativePath() = %q", got)
		}
		if got := v.ThumbnailRelativePath(); got != filepath.Join("/clips", id, "thumb.jpg") {
			t.Errorf("ThumbnailRelativePath() = %q", got)
		}
		if got := v.ThumbnailXSRelativePath(); got != filepath.Join("/clips", id, "thumb-xs.jpg") {
			t.Errorf("ThumbnailXSRelativePath() = %q", got)
		}
	})
	t.Run("video", func(t *testing.T) {
		v := &Video{ID: id, Type: "v"}
		if got := v.FolderRelativePath(); got != filepath.Join("/videos", id) {
			t.Errorf("FolderRelativePath() = %q", got)
		}
	})
	t.Run("movie", func(t *testing.T) {
		v := &Video{ID: id, Type: "m"}
		if got := v.FolderRelativePath(); got != filepath.Join("/movies", id) {
			t.Errorf("FolderRelativePath() = %q", got)
		}
	})
}
