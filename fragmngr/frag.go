package fragmngr

import (
	"errors"
	"sync"
)

var (
	//ErrFragIndexOutOfRange -- fragment index out of range
	ErrFragIndexOutOfRange = errors.New("fragment index out of range")
	//ErrFragNotExist -- fragment data not exist
	ErrFragNotExist = errors.New("fragment data not exist")
	//ErrInvalidParameter -- parameter invalid
	ErrInvalidParameter = errors.New("parameter is invalid")
	//ErrWriteDBFail -- failed to write to db
	ErrWriteDBFail = errors.New("failed to write to db")
)

//FileAdapter operate a file to store and provide fragment data
type FileAdapter interface {
	//Read reads out fragment data at the give position
	Read(index uint64) ([]byte, error)
	//Write stores data to the given position
	Write(index uint64, d []byte) error

	//Exists checks whether the data of given postion exists already in local
	Exists(index uint64) (bool, error)
	//FragCount returns total fragment count
	FragCount() uint64

	//Close closes the underlying opened file and save all changes
	Close() error
}

//FragManager is based on underlying mechanism and returns FileAdapter of given file name
type FragManager interface {
	//Open returns the FileAdapter for futher operation
	//the file will be created if it doesnt exist
	Open(name string) (FileAdapter, error)
	//Remove removes the underlying file
	Remove(name string) error
}

//====================================================================
var holder fragManagerHolder

type fragManagerHolder struct {
	lock        sync.Mutex
	fragManager FragManager
}

func (fmh *fragManagerHolder) setFragManager(fm FragManager) bool {
	fmh.lock.Lock()
	defer fmh.lock.Unlock()

	if fmh.fragManager != nil {
		return false
	}
	fmh.fragManager = fm
	return true
}

func (fmh *fragManagerHolder) Instance() FragManager {
	return fmh.fragManager
}

//FragManagerInstance returns the current FragManager instance
func FragManagerInstance() FragManager {
	return holder.Instance()
}
