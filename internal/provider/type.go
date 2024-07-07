package provider

type Provider interface {
	// actors
	ActorSearch(string) (string, error)
	ActorGetThumb(string, string) ([]byte, error)

	// slug for registering
	SlugGet() string

	// nice name for frontend
	NiceName() string
}
