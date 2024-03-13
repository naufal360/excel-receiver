package util

import "os"

// CreateDirectory create multiple directory.
func CreateDirectory(paths ...string) (err error) {
	for _, path := range paths {
		_, notExistError := os.Stat(path)
		if os.IsNotExist(notExistError) {
			if err = os.MkdirAll(path, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return
}
