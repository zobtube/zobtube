package storage

import (
	"io"
	"path/filepath"
)

// CopyObject copies an object from src at srcPath to dst at dstPath.
func CopyObject(src, dst Storage, srcPath, dstPath string) error {
	if err := dst.MkdirAll(filepath.Dir(dstPath)); err != nil {
		return err
	}
	rc, err := src.Open(srcPath)
	if err != nil {
		return err
	}
	defer rc.Close()
	wc, err := dst.Create(dstPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(wc, rc)
	if cerr := wc.Close(); err == nil {
		err = cerr
	}
	return err
}
