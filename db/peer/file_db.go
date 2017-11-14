package peer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"

	"github.com/Zumium/fyer/merkle"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	filesPrefix = "files."
)

var (
	//ErrDecodingFail happens when database bytes cannot be decoded
	ErrDecodingFail = errors.New("failed to decode database bytes")
)

//FilesDBWrapper provides convenient operations on database file records
type FilesDBWrapper struct {
	Name string

	db *leveldb.DB
}

//NewFilesDBWrapper creates a lightweight db wrapper to simplify db operations
func NewFilesDBWrapper(name string) *FilesDBWrapper {
	return &FilesDBWrapper{Name: name, db: instance()}
}

//Has returns true if the file record exists
func (fw *FilesDBWrapper) Has() (bool, error) {
	ret, err := fw.db.Has(fw.key(""), nil)
	if err != nil {
		return false, err
	}
	return ret, nil
}

//Size -- file size
func (fw *FilesDBWrapper) Size() (uint64, error) {
	val, err := fw.db.Get(fw.key("size"), nil)
	if err != nil {
		return 0, err
	}
	size, n := binary.Uvarint(val)
	if n <= 0 {
		return 0, ErrDecodingFail
	}
	return size, nil
}

//Hash -- file hash
func (fw *FilesDBWrapper) Hash() ([]byte, error) {
	val, err := fw.db.Get(fw.key("hash"), nil)
	if err != nil {
		return nil, err
	}
	return val, nil
}

//FragCount -- total number of fragments
func (fw *FilesDBWrapper) FragCount() (uint64, error) {
	val, err := fw.db.Get(fw.key("frag_count"), nil)
	if err != nil {
		return 0, err
	}
	fragCount, n := binary.Uvarint(val)
	if n <= 0 {
		return 0, ErrDecodingFail
	}
	return fragCount, nil
}

//UploadTime -- file uploading time
func (fw *FilesDBWrapper) UploadTime() (time.Time, error) {
	var t time.Time
	val, err := fw.db.Get(fw.key("upload_time"), nil)
	if err != nil {
		return t, err
	}
	if err := t.UnmarshalBinary(val); err != nil {
		return t, err
	}
	return t, nil
}

//MerkleTree -- file merkle tree
func (fw *FilesDBWrapper) MerkleTree() (*merkle.MTree, error) {
	val, err := fw.db.Get(fw.key("merkle_tree"), nil)
	if err != nil {
		return nil, err
	}
	mtree, err := merkle.Unmarshal(val)
	if err != nil {
		return nil, err
	}
	return mtree, nil
}

//StoredFrags -- file's stored fragment numbers
func (fw *FilesDBWrapper) StoredFrags() (*StoredFrags, error) {
	val, err := fw.db.Get(fw.key("stored_frags"), nil)
	if err != nil {
		return nil, err
	}
	sf := NewEmptyStoredFrags()
	if err := sf.UnmarshalJSON(val); err != nil {
		return nil, err
	}
	return sf, nil
}

//Edit opens an underlying transaction to allow modifying the database record
func (fw *FilesDBWrapper) Edit() *FilesDBEditor {
	return &FilesDBEditor{
		wrapper: fw,
		batch:   new(leveldb.Batch),
	}
}

func (fw *FilesDBWrapper) key(keyname string) []byte {
	var (
		err    error
		keyBuf bytes.Buffer
	)
	_, err = keyBuf.WriteString(filesPrefix)
	if err != nil {
		panic(err)
	}
	_, err = keyBuf.WriteString(fw.Name)
	if err != nil {
		panic(err)
	}
	if keyname != "" {
		_, err = keyBuf.WriteString(".")
		if err != nil {
			panic(err)
		}
		_, err = keyBuf.WriteString(keyname)
		if err != nil {
			panic(err)
		}
	}
	return keyBuf.Bytes()
}

//===============================================================

//FilesDBEditor is a helper to modifying database record
type FilesDBEditor struct {
	wrapper *FilesDBWrapper
	batch   *leveldb.Batch

	err error
}

//Done shall be called when setting operations are done
func (editor *FilesDBEditor) Done() error {
	if editor.err != nil {
		return editor.err
	}
	editor.batch.Put(editor.wrapper.key(""), []byte{1})
	return editor.wrapper.db.Write(editor.batch, nil)
}

//SetSize stores file size
func (editor *FilesDBEditor) SetSize(size uint64) *FilesDBEditor {
	if editor.err != nil {
		return editor
	}

	val := make([]byte, 8)
	if n := binary.PutUvarint(val, size); n <= 0 {
		panic("internal error: failed to write uint64 to []byte")
	}
	editor.batch.Put(editor.wrapper.key("size"), val)
	return editor
}

//SetHash stores file hash
func (editor *FilesDBEditor) SetHash(hash []byte) *FilesDBEditor {
	if editor.err != nil {
		return editor
	}

	editor.batch.Put(editor.wrapper.key("hash"), hash)
	return editor
}

//SetFragCount stores total number of fragments
func (editor *FilesDBEditor) SetFragCount(fragCount uint64) *FilesDBEditor {
	if editor.err != nil {
		return editor
	}

	val := make([]byte, 8)
	if n := binary.PutUvarint(val, fragCount); n <= 0 {
		panic("internal error: failed to write uint64 to []byte")
	}
	editor.batch.Put(editor.wrapper.key("frag_count"), val)
	return editor
}

//SetUploadTime stores uploading time
func (editor *FilesDBEditor) SetUploadTime(uploadTime time.Time) *FilesDBEditor {
	if editor.err != nil {
		return editor
	}

	val, err := uploadTime.MarshalBinary()
	if err != nil {
		editor.err = err
		return editor
	}
	editor.batch.Put(editor.wrapper.key("upload_time"), val)
	return editor
}

//SetMerkleTree stores the file's merkle tree
func (editor *FilesDBEditor) SetMerkleTree(mtree *merkle.MTree) *FilesDBEditor {
	if editor.err != nil {
		return editor
	}

	val, err := merkle.Marshal(mtree)
	if err != nil {
		editor.err = err
		return editor
	}
	editor.batch.Put(editor.wrapper.key("merkle_tree"), val)
	return editor
}

//SetStoredFrags stores numbers of stored fragments
func (editor *FilesDBEditor) SetStoredFrags(sf *StoredFrags) *FilesDBEditor {
	if editor.err != nil {
		return editor
	}

	val, err := sf.MarshalJSON()
	if err != nil {
		editor.err = err
		return editor
	}
	editor.batch.Put(editor.wrapper.key("stored_frags"), val)
	return editor
}

//===============================================================

//NewFileRecord creates a new file database record
func NewFileRecord(name string, size uint64, hash []byte, fragCount uint64, uploadTime time.Time, mtree *merkle.MTree, storedFrags *StoredFrags) (*FilesDBWrapper, error) {
	nfdw := NewFilesDBWrapper(name)
	err := nfdw.Edit().SetSize(size).SetHash(hash).SetFragCount(fragCount).SetUploadTime(uploadTime).SetMerkleTree(mtree).SetStoredFrags(storedFrags).Done()
	return nfdw, err
}
