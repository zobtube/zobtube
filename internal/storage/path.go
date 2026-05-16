package storage

import "path/filepath"

// FirstExistingPath returns the first path in paths that exists on store.
func FirstExistingPath(store Storage, paths []string) (string, bool, error) {
	for _, p := range paths {
		ok, err := store.Exists(p)
		if err != nil {
			return "", false, err
		}
		if ok {
			return p, true, nil
		}
	}
	return "", false, nil
}

// WalkFiles lists files under prefix. When recursive is true, subdirectories are traversed.
// fn receives the storage-relative path (e.g. triage/foo/bar.mp4) and the entry metadata.
func WalkFiles(s Storage, prefix string, recursive bool, fn func(path string, e Entry) error) error {
	entries, err := s.List(prefix)
	if err != nil {
		return err
	}
	for _, e := range entries {
		p := filepath.Join(prefix, e.Name)
		if e.IsDir {
			if !recursive {
				continue
			}
			if err := WalkFiles(s, p, true, fn); err != nil {
				return err
			}
			continue
		}
		if err := fn(p, e); err != nil {
			return err
		}
	}
	return nil
}
