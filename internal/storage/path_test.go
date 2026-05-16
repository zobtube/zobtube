package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWalkFiles_Recursive(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "triage", "sub"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "triage", "a.mp4"), []byte("a"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "triage", "sub", "b.mp4"), []byte("b"), 0o644)

	fs := NewFilesystem(root)
	var paths []string
	err := WalkFiles(fs, "triage", true, func(p string, e Entry) error {
		if !e.IsDir {
			paths = append(paths, p)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 {
		t.Fatalf("paths = %v, want 2 files", paths)
	}

	paths = nil
	err = WalkFiles(fs, "triage", false, func(p string, e Entry) error {
		if !e.IsDir {
			paths = append(paths, p)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 1 {
		t.Fatalf("non-recursive paths = %v, want 1", paths)
	}
}
