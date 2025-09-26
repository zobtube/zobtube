package provider

import (
	"errors"
)

type Provider interface {
	// actors
	ActorSearch(bool, string) (string, error)
	ActorGetThumb(bool, string, string) ([]byte, error)

	// slug for registering
	SlugGet() string

	// nice name for frontend
	NiceName() string

	// Capabilities
	CapabilitySearchActor() bool
	CapabilityScrapePicture() bool
}

var ErrOfflineMode = errors.New("not possible in offline mode")
