package storage

import (
	"io"
	"time"
)

// Entry describes a single item from List.
type Entry struct {
	Name  string
	Size  int64
	ModTime time.Time
	IsDir bool
}

// Storage is the abstraction for library media storage (filesystem or S3).
type Storage interface {
	// Open opens the object at path for reading. Caller must close the returned ReadCloser.
	Open(path string) (io.ReadCloser, error)
	// Create creates or overwrites the object at path. Caller must close the returned WriteCloser.
	Create(path string) (io.WriteCloser, error)
	// Delete removes the object at path. No error if it does not exist.
	Delete(path string) error
	// Exists returns true if the object at path exists.
	Exists(path string) (bool, error)
	// MkdirAll ensures the prefix path exists (no-op for S3).
	MkdirAll(path string) error
	// List returns entries under prefix. Paths are relative to the storage root.
	List(prefix string) ([]Entry, error)
}
