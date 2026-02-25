package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupActorAliasDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&ActorAlias{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestActorAlias_BeforeCreate_AssignsUUIDWhenIDEmpty(t *testing.T) {
	db := setupActorAliasDB(t)
	a := &ActorAlias{Name: "alias", ActorID: "b23f4f4a-1c5c-11f0-8822-305a3a05e04d"}
	if err := db.Create(a).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if a.ID == "" {
		t.Error("expected ID to be set by BeforeCreate")
	}
	if _, err := uuid.Parse(a.ID); err != nil {
		t.Errorf("expected ID to be valid UUID, got %q: %v", a.ID, err)
	}
}

func TestActorAlias_BeforeCreate_AssignsUUIDWhenZeroUUID(t *testing.T) {
	db := setupActorAliasDB(t)
	a := &ActorAlias{ID: "00000000-0000-0000-0000-000000000000", Name: "alias", ActorID: "b23f4f4a-1c5c-11f0-8822-305a3a05e04d"}
	if err := db.Create(a).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if a.ID == "00000000-0000-0000-0000-000000000000" {
		t.Error("expected zero UUID to be replaced")
	}
	if _, err := uuid.Parse(a.ID); err != nil {
		t.Errorf("expected ID to be valid UUID, got %q: %v", a.ID, err)
	}
}

func TestActorAlias_BeforeCreate_KeepsIDWhenSet(t *testing.T) {
	db := setupActorAliasDB(t)
	wantID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	a := &ActorAlias{ID: wantID, Name: "alias", ActorID: "b23f4f4a-1c5c-11f0-8822-305a3a05e04d"}
	if err := db.Create(a).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if a.ID != wantID {
		t.Errorf("expected ID to remain %q, got %q", wantID, a.ID)
	}
}
