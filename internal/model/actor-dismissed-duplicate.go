package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ActorDismissedDuplicate records an admin decision that two actors sharing a name
// are not duplicates and should not be suggested again.
type ActorDismissedDuplicate struct {
	ID        string `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	ActorID1  string `gorm:"type:uuid;index"`
	ActorID2  string `gorm:"type:uuid;index"`
}

// NormalizeActorPair returns actor IDs in lexicographic order for stable storage.
func NormalizeActorPair(id1, id2 string) (string, string) {
	if id1 < id2 {
		return id1, id2
	}
	return id2, id1
}

func (a *ActorDismissedDuplicate) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "00000000-0000-0000-0000-000000000000" {
		a.ID = uuid.NewString()
		return nil
	}
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	return nil
}
