package storage

import (
	"io"
	"os"
	"path/filepath"
)

// Filesystem implements Storage using a local directory root.
type Filesystem struct {
	Root string
}

// NewFilesystem returns a Storage that uses the given root directory.
func NewFilesystem(root string) *Filesystem {
	return &Filesystem{Root: root}
}

// Open opens the file at path for reading.
func (f *Filesystem) Open(path string) (io.ReadCloser, error) {
	full := filepath.Join(f.Root, path)
	return os.Open(full)
}

// Create creates or overwrites the file at path.
func (f *Filesystem) Create(path string) (io.WriteCloser, error) {
	full := filepath.Join(f.Root, path)
	dir := filepath.Dir(full)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, err
	}
	return os.Create(full)
}

// Delete removes the file at path.
func (f *Filesystem) Delete(path string) error {
	full := filepath.Join(f.Root, path)
	err := os.Remove(full)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return err
}

// Exists returns true if the file at path exists.
func (f *Filesystem) Exists(path string) (bool, error) {
	full := filepath.Join(f.Root, path)
	_, err := os.Stat(full)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// MkdirAll ensures the directory at path exists.
func (f *Filesystem) MkdirAll(path string) error {
	full := filepath.Join(f.Root, path)
	return os.MkdirAll(full, 0o750)
}

// List returns entries under prefix. prefix is relative to Root.
func (f *Filesystem) List(prefix string) ([]Entry, error) {
	full := filepath.Join(f.Root, prefix)
	entries, err := os.ReadDir(full)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		out = append(out, Entry{
			Name:   e.Name(),
			Size:   info.Size(),
			ModTime: info.ModTime(),
			IsDir:  e.IsDir(),
		})
	}
	return out, nil
}

// FullPath returns the absolute path for a relative path (for use with gin.Context.File when serving).
// Only valid for filesystem storage.
func (f *Filesystem) FullPath(path string) string {
	return filepath.Join(f.Root, path)
}

// Stat returns FileInfo for the path if the storage is filesystem-based.
// Returns nil if not available (e.g. S3). Used by controllers that serve via g.File().
func (f *Filesystem) Stat(path string) (os.FileInfo, error) {
	full := filepath.Join(f.Root, path)
	return os.Stat(full)
}
