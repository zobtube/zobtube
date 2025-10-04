package config

import (
	"os"
	"path/filepath"
)

func (cfg *Config) EnsureTreePresent() error {
	folders := []string{
		"clips",
		"movies",
		"videos",
		"actors",
		"triage",
	}

	// ensure library folder exists
	path := cfg.Media.Path
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// do not exists, create it
		err = os.Mkdir(path, 0o750)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// ensure folders inside the library exist
	for _, folder := range folders {
		path := filepath.Join(cfg.Media.Path, folder)
		// ensure folder exists
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			// do not exists, create it
			err = os.Mkdir(path, 0o750)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}
