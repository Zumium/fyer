package filemanager

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

var (
	ErrNotFound = errors.New("file not found in registration")
)

//files contains the map of file name to absolute file path
var files = make(map[string]string)

//Register registers a file
func Register(fpath string) error {
	absPath, err := filepath.Abs(fpath)
	if err != nil {
		return err
	}
	files[filepath.Base(absPath)] = absPath
	return nil
}

//Unregister deletes the record
func Unregister(name string) bool {
	_, exist := files[name]
	if !exist {
		return false
	}
	delete(files, name)
	return true
}

//Open opens the specified file
func Open(name string) (*os.File, error) {
	fpath, exist := files[name]
	if !exist {
		return nil, ErrNotFound
	}
	return os.Open(fpath)
}
