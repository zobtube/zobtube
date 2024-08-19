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

	for _, folder := range folders {
		path := filepath.Join(cfg.Media.Path, folder)
		// ensure folder exists
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			// do not exists, create it
			err = os.Mkdir(path, os.ModePerm)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}
