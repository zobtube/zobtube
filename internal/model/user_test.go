package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupUserDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestUser_BeforeCreate_AssignsUUIDWhenIDEmpty(t *testing.T) {
	db := setupUserDB(t)
	u := &User{Username: "u", Password: "p"}
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

func TestUser_BeforeCreate_AssignsUUIDWhenZeroUUID(t *testing.T) {
	db := setupUserDB(t)
	u := &User{ID: "00000000-0000-0000-0000-000000000000", Username: "u", Password: "p"}
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

func TestUser_BeforeCreate_KeepsIDWhenSet(t *testing.T) {
	db := setupUserDB(t)
	wantID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	u := &User{ID: wantID, Username: "u", Password: "p"}
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if u.ID != wantID {
		t.Errorf("expected ID to remain %q, got %q", wantID, u.ID)
	}
}

func TestUser_URLAdmDelete(t *testing.T) {
	u := &User{ID: "deadbeef-1234-5678-abcd-000000000000"}
	got := u.URLAdmDelete()
	want := "/adm/user/deadbeef-1234-5678-abcd-000000000000/delete"
	if got != want {
		t.Errorf("URLAdmDelete() = %q, want %q", got, want)
	}
}
