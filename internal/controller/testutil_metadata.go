package controller

import (
	"github.com/zobtube/zobtube/internal/storage"
)

// registerTestMetadataStorage attaches a filesystem metadata store for tests.
func registerTestMetadataStorage(ctrl *Controller, root string) {
	ctrl.MetadataStorageRegister(storage.NewFilesystem(root))
}
