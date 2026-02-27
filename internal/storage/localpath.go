package storage

import (
	"io"
	"os"
	"path/filepath"
)

// LocalPathForRead returns a local filesystem path that can be used to read the object (e.g. by ffprobe/ffmpeg).
// For filesystem storage this is the actual path; for S3 etc. the object is copied to a temp file.
// The caller must call the returned cleanup function when done.
func LocalPathForRead(store Storage, path string) (localPath string, cleanup func(), err error) {
	if store == nil {
		return "", nil, os.ErrInvalid
	}
	if fs, ok := store.(*Filesystem); ok {
		return fs.FullPath(path), func() {}, nil
	}
	rc, err := store.Open(path)
	if err != nil {
		return "", nil, err
	}
	defer rc.Close()
	f, err := os.CreateTemp("", "zt-*"+filepath.Base(path))
	if err != nil {
		return "", nil, err
	}
	_, err = io.Copy(f, rc)
	if err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", nil, err
	}
	if err := f.Close(); err != nil {
		os.Remove(f.Name())
		return "", nil, err
	}
	return f.Name(), func() { os.Remove(f.Name()) }, nil
}
