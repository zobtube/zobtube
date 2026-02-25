package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTaskDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&Task{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestTask_BeforeCreate_AssignsUUIDWhenIDEmpty(t *testing.T) {
	db := setupTaskDB(t)
	tt := &Task{Name: "task", Step: "step"}
	if err := db.Create(tt).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if tt.ID == "" {
		t.Error("expected ID to be set by BeforeCreate")
	}
	if _, err := uuid.Parse(tt.ID); err != nil {
		t.Errorf("expected ID to be valid UUID, got %q: %v", tt.ID, err)
	}
}

func TestTask_BeforeCreate_AssignsUUIDWhenZeroUUID(t *testing.T) {
	db := setupTaskDB(t)
	tt := &Task{ID: "00000000-0000-0000-0000-000000000000", Name: "task", Step: "step"}
	if err := db.Create(tt).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if tt.ID == "00000000-0000-0000-0000-000000000000" {
		t.Error("expected zero UUID to be replaced")
	}
	if _, err := uuid.Parse(tt.ID); err != nil {
		t.Errorf("expected ID to be valid UUID, got %q: %v", tt.ID, err)
	}
}

func TestTask_BeforeCreate_KeepsIDWhenSet(t *testing.T) {
	db := setupTaskDB(t)
	wantID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	tt := &Task{ID: wantID, Name: "task", Step: "step"}
	if err := db.Create(tt).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if tt.ID != wantID {
		t.Errorf("expected ID to remain %q, got %q", wantID, tt.ID)
	}
}
