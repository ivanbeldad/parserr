package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindFile Search for a file and return either its location or an error
func FindFile(root, filename string) (location string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.Name() == filename {
			location = path
			return fmt.Errorf("ok")
		}
		return nil
	})
	if err != nil && err.Error() == "ok" {
		err = nil
	}
	if location == "" {
		err = fmt.Errorf("%s doesn't exists inside %s", filename, root)
	}
	return
}
