package fragmngr

import (
	"os"
	"path/filepath"

	peer_db "github.com/Zumium/fyer/db/peer"
)

type simpleFSFragManager struct {
	base string
}

//InitSimpleFSFragManager creates and use simple filesystem fragment manager
func InitSimpleFSFragManager(basePath string) error {
	holder.setFragManager(&simpleFSFragManager{base: basePath})
	return nil //NO NEED TO RETURN ANY ERRORS
}

//Open opens (or creates) a file on the filesystem
func (m *simpleFSFragManager) Open(name string) (FileAdapter, error) {
	return newSimpleFSFileAdapter(base, name)
}

//Remove removes the underlying file from fs
func (m *simpleFSFragManager) Remove(name string) error {
	return os.Remove(filepath.Join(base, name))
}

//---------------------------------------------------------

type simpleFSFileAdapter struct {
	fragSize uint32
	file     *os.File
}

func newSimpleFSFileAdapter(base string, name string) (*simpleFSFileAdapter, error) {
	f, err := os.Open(filepath.Join(base, name))
	if err != nil {
		return nil, err
	}
	return nil, &simpleFSFileAdapter{
		file:     f,
		fragSize: 2 * 1024 * 1024,
	}
}

func (ada *simpleFSFileAdapter) fileDBWrapper() *peer_db.FilesDBWrapper {
	return &peer_db.NewFilesDBWrapper(filepath.Base(ada.file.Name()))
}

func (ada *simpleFSFileAdapter) checkIndex(index uint64) (valid, last bool) {
	if diff := ada.fragCount() - 1 - index; diff > 0 {
		valid, last = false, false
	} else if diff == 0 {
		valid, last = true, true
	} else {
		valid, last = true, false
	}
	return
}

func (ada *simpleFSFileAdapter) size() uint64 {
	//QUERYING FILE SIZE FROM DATABASE
	size, err := ada.fileDBWrapper().Size()
	if err != nil {
		panic(err) //Maybe not so good
	}
	return size
}

func (ada *simpleFSFileAdapter) fragCount() uint64 {
	//QUERYING FRAGMENT COUNT FROM DATABASE
	fcount, err := ada.fileDBWrapper().FragCount()
	if err != nil {
		panic(err) //Maybe not so good
	}
	return fcount
}

func (ada *simpleFSFileAdapter) Read(index uint64) (d []byte, err error) {
	valid, last := ada.checkIndex(index)
	if !valid {
		return nil, ErrFragIndexOutOfRange
	}

	//CHECK WHETHER THE FRAGMENT EXISTS IN LOCAL IN DATABASE
	var storedFrags *peer_db.StoredFrags
	wrapper := ada.fileDBWrapper()
	storedFrags, err = wrapper.StoredFrags()
	if err != nil {
		return
	}
	if !storedFrags.Has() {
		d, err = nil, ErrFragNotExist
		return
	}

	begin := index * ada.fragSize
	if last {
		d = make([]byte, ada.size()-begin)
	} else {
		d = make([]byte, ada.fragSize)
	}

	if _, err = ada.file.ReadAt(d, begin); err != nil {
		d = nil
		return
	}
	err = nil
	return
}

func (ada *simpleFSFileAdapter) Write(index uint64, d []byte) error {
	valid, last := ada.checkIndex(index)
	if !valid {
		return ErrFragIndexOutOfRange
	}

	begin := index * ada.fragSize
	s := uint64(0)
	if last {
		s = ada.size() - begin
	} else {
		s = ada.fragSize
	}
	if len(d) != s {
		return
	}

	if _, err := ada.file.WriteAt(d, begin); err != nil {
		return err
	}

	//ADD RECORD OF THE NEWLY STORED FRAGMENT IN DATABASE
	wrapper := ada.fileDBWrapper()
	storedFrags, err := wrapper.StoredFrags()
	if err != nil {
		return err
	}
	storedFrags.Add(index)
	return wrapper.Edit().SetStoredFrags(storedFrags).Done()
}

func (ada *simpleFSFileAdapter) Exists(index uint64) (bool, error) {
	//QUERYING STORED FRAGMENT NUMBERS FROM DATABASE
	storedFrags, err := ada.fileDBWrapper().StoredFrags()
	if err != nil {
		return false, err
	}
	return storedFrags.Has(index), nil
}

func (ada *simpleFSFileAdapter) FragCount() uint64 {
	return ada.fragCount()
}

func (ada *simpleFSFileAdapter) Close() error {
	return ada.file.Close()
}
