package storage

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
