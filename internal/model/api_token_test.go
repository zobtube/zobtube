package model

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupApiTokenDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&ApiToken{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestApiToken_BeforeCreate_AssignsUUIDWhenIDEmpty(t *testing.T) {
	db := setupApiTokenDB(t)
	token := &ApiToken{
		ID:        "",
		UserID:    "b23f4f4a-1c5c-11f0-8822-305a3a05e04d",
		Name:      "test",
		TokenHash: "abc123",
		CreatedAt: time.Now(),
	}
	if err := db.Create(token).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if token.ID == "" {
		t.Error("expected ID to be set by BeforeCreate")
	}
	if _, err := uuid.Parse(token.ID); err != nil {
		t.Errorf("expected ID to be valid UUID, got %q: %v", token.ID, err)
	}
}

func TestApiToken_BeforeCreate_KeepsIDWhenSet(t *testing.T) {
	db := setupApiTokenDB(t)
	wantID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	token := &ApiToken{
		ID:        wantID,
		UserID:    "b23f4f4a-1c5c-11f0-8822-305a3a05e04d",
		Name:      "test",
		TokenHash: "def456",
		CreatedAt: time.Now(),
	}
	if err := db.Create(token).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if token.ID != wantID {
		t.Errorf("expected ID to remain %q, got %q", wantID, token.ID)
	}
}
