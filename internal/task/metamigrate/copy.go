package metamigrate

import (
	"io"
	"os"
	"path/filepath"

	"github.com/zobtube/zobtube/internal/storage"
)

func copyStorageObject(src, dst storage.Storage, path string) error {
	if err := dst.MkdirAll(filepath.Dir(path)); err != nil {
		return err
	}
	rc, err := src.Open(path)
	if err != nil {
		return err
	}
	defer rc.Close()
	wc, err := dst.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(wc, rc)
	if cerr := wc.Close(); err == nil {
		err = cerr
	}
	return err
}

func removeStorageObject(store storage.Storage, path string) error {
	return store.Delete(path)
}

// removeEmptyDirFilesystem removes relDir when empty (filesystem backends only).
func removeEmptyDirFilesystem(store storage.Storage, relDir string) {
	fs, ok := store.(*storage.Filesystem)
	if !ok || relDir == "" || relDir == "." {
		return
	}
	entries, err := fs.List(relDir)
	if err != nil || len(entries) > 0 {
		return
	}
	full := fs.FullPath(relDir)
	if err := os.Remove(full); err != nil {
		return
	}
	parent := filepath.Dir(relDir)
	if parent != "." && parent != "/" && parent != relDir {
		removeEmptyDirFilesystem(store, parent)
	}
}

// cleanupLegacyEntityRoots removes top-level entity thumbnail folders when empty (filesystem only).
func cleanupLegacyEntityRoots(store storage.Storage) {
	for _, dir := range []string{"actors", "channels", "categories"} {
		removeEmptyDirFilesystem(store, dir)
	}
}
