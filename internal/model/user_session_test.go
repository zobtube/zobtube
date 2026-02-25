package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupUserSessionDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&UserSession{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestUserSession_BeforeCreate_AssignsUUIDWhenIDEmpty(t *testing.T) {
	db := setupUserSessionDB(t)
	u := &UserSession{ValidUntil: time.Now()}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if u.ID == "" {
		t.Error("expected ID to be set by BeforeCreate")
	}
	if _, err := uuid.Parse(u.ID); err != nil {
		t.Errorf("expected ID to be valid UUID, got %q: %v", u.ID, err)
	}
}

func TestUserSession_BeforeCreate_AssignsUUIDWhenZeroUUID(t *testing.T) {
	db := setupUserSessionDB(t)
	u := &UserSession{ID: "00000000-0000-0000-0000-000000000000", ValidUntil: time.Now()}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if u.ID == "00000000-0000-0000-0000-000000000000" {
		t.Error("expected zero UUID to be replaced")
	}
	if _, err := uuid.Parse(u.ID); err != nil {
		t.Errorf("expected ID to be valid UUID, got %q: %v", u.ID, err)
	}
}

func TestUserSession_BeforeCreate_KeepsIDWhenSet(t *testing.T) {
	db := setupUserSessionDB(t)
	wantID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	u := &UserSession{ID: wantID, ValidUntil: time.Now()}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if u.ID != wantID {
		t.Errorf("expected ID to remain %q, got %q", wantID, u.ID)
	}
}
