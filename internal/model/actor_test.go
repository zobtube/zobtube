package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupActorDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&Actor{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestActor_BeforeCreate_AssignsUUIDWhenIDEmpty(t *testing.T) {
	db := setupActorDB(t)
	a := &Actor{Name: "actor"}
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

func TestActor_BeforeCreate_AssignsUUIDWhenZeroUUID(t *testing.T) {
	db := setupActorDB(t)
	a := &Actor{ID: "00000000-0000-0000-0000-000000000000", Name: "actor"}
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

func TestActor_BeforeCreate_KeepsIDWhenSet(t *testing.T) {
	db := setupActorDB(t)
	wantID := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
	a := &Actor{ID: wantID, Name: "actor"}
	if err := db.Create(a).Error; err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if a.ID != wantID {
		t.Errorf("expected ID to remain %q, got %q", wantID, a.ID)
	}
}

func TestActor_SexTypeAsString(t *testing.T) {
	tests := []struct {
		sex  string
		want string
	}{
		{"m", "male"},
		{"f", "female"},
		{"tw", "trans-women"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.sex, func(t *testing.T) {
			a := &Actor{Sex: tt.sex}
			got := a.SexTypeAsString()
			if got != tt.want {
				t.Errorf("SexTypeAsString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestActor_AliasesAsNiceString(t *testing.T) {
	tests := []struct {
		name   string
		aliases []ActorAlias
		want   string
	}{
		{"empty", nil, ""},
		{"one", []ActorAlias{{Name: "Alias1"}}, "Alias1"},
		{"many", []ActorAlias{{Name: "A"}, {Name: "B"}, {Name: "C"}}, "A / B / C"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Actor{Aliases: tt.aliases}
			got := a.AliasesAsNiceString()
			if got != tt.want {
				t.Errorf("AliasesAsNiceString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestActor_URLView_URLThumb_URLAdmEdit_URLAdmDelete(t *testing.T) {
	id := "deadbeef-1234-5678-abcd-000000000000"
	a := &Actor{ID: id}
	if got := a.URLView(); got != "/actor/"+id {
		t.Errorf("URLView() = %q", got)
	}
	if got := a.URLThumb(); got != "/api/actor/"+id+"/thumb" {
		t.Errorf("URLThumb() = %q", got)
	}
	if got := a.URLAdmEdit(); got != "/actor/"+id+"/edit" {
		t.Errorf("URLAdmEdit() = %q", got)
	}
	if got := a.URLAdmDelete(); got != "/actor/"+id+"/delete" {
		t.Errorf("URLAdmDelete() = %q", got)
	}
}
