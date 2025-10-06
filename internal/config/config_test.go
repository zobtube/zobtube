package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog/log"
)

// --- Tests for New() ---

func TestNew_Success(t *testing.T) {
	logger := log.Logger

	cfg, err := New(&logger, ":8080", "sqlite", "file:test.db", "/tmp")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Server.Bind != ":8080" {
		t.Errorf("expected Bind to be ':8080', got %q", cfg.Server.Bind)
	}
	if cfg.DB.Driver != "sqlite" {
		t.Errorf("expected DB.Driver to be 'sqlite', got %q", cfg.DB.Driver)
	}
	if cfg.DB.Connstring != "file:test.db" {
		t.Errorf("expected DB.Connstring to be 'file:test.db', got %q", cfg.DB.Connstring)
	}
	if cfg.Media.Path != "/tmp" {
		t.Errorf("expected Media.Path to be '/tmp', got %q", cfg.Media.Path)
	}
}

func TestNew_MissingFields(t *testing.T) {
	logger := log.Logger
	tests := []struct {
		name       string
		driver     string
		connstring string
		path       string
		wantErr    string
	}{
		{"missing driver", "", "conn", "/tmp", "ZT_DB_DRIVER is not set"},
		{"missing connstring", "sqlite", "", "/tmp", "ZT_DB_CONNSTRING is not set"},
		{"missing media path", "sqlite", "conn", "", "ZT_MEDIA_PATH is not set"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(&logger, ":8080", tt.driver, tt.connstring, tt.path)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !errors.Is(err, errors.New(tt.wantErr)) && err.Error() != tt.wantErr {
				t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

// --- Tests for EnsureTreePresent() ---

func TestEnsureTreePresent_CreatesMissingFolders(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{}
	cfg.Media.Path = filepath.Join(tmpDir, "media")

	err := cfg.EnsureTreePresent()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check main folder
	if _, err := os.Stat(cfg.Media.Path); os.IsNotExist(err) {
		t.Error("expected media path to exist, got not exist")
	}

	// Check subfolders
	subfolders := []string{"clips", "movies", "videos", "actors", "triage"}
	for _, folder := range subfolders {
		path := filepath.Join(cfg.Media.Path, folder)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected folder %s to exist", folder)
		}
	}
}

func TestEnsureTreePresent_ExistingFolders(t *testing.T) {
	tmpDir := t.TempDir()
	mediaPath := filepath.Join(tmpDir, "media")

	// Pre-create structure
	if err := os.Mkdir(mediaPath, 0o750); err != nil {
		t.Fatalf("failed to create media folder: %v", err)
	}
	subfolders := []string{"clips", "movies", "videos", "actors", "triage"}
	for _, f := range subfolders {
		err := os.Mkdir(filepath.Join(mediaPath, f), 0o750)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}

	cfg := &Config{}
	cfg.Media.Path = mediaPath

	err := cfg.EnsureTreePresent()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestEnsureTreePresent_InvalidPath(t *testing.T) {
	// Simulate a case where os.Mkdir fails due to invalid path
	cfg := &Config{}
	cfg.Media.Path = "/this/should/not/exist/and/cannot/be/created"

	err := cfg.EnsureTreePresent()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
