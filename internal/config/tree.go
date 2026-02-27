package config

import (
	"os"
	"path/filepath"
)

var defaultLibraryFolders = []string{
	"clips",
	"movies",
	"videos",
	"actors",
	"triage",
}

// EnsureTreePresent ensures the library folder and default subfolders exist at path.
// Used for the single configured media path (backward compat) or for each filesystem library path.
func (cfg *Config) EnsureTreePresent() error {
	return EnsureTreePresentForPath(cfg.Media.Path)
}

// EnsureTreePresentForPath ensures the library folder and default subfolders exist at the given path.
func EnsureTreePresentForPath(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0o750)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	for _, folder := range defaultLibraryFolders {
		dir := filepath.Join(path, folder)
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, 0o750)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}
