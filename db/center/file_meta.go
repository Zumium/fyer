package center

import (
	"time"

	"github.com/Zumium/fyer/merkle"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//fileMetaC returns the file meta collection
func fileMetaC() *mgo.Collection {
	return mgoFyerDB().C(mgoFileMetaCollection)
}

type fileMetaRecordMode uint8

const (
	fileMetaModeNew fileMetaRecordMode = iota
	fileMetaModeNormal
)

//FileMeta represents a database record of file meta data
type FileMeta struct {
	name string

	mode fileMetaRecordMode
	doc  mgoFileMeta
	err  error
}

//ToFileMeta creates a new FileMeta to apply furthur db operaions
func ToFileMeta(name string) (fmeta *FileMeta) {
	fmeta = &FileMeta{name: name}
	fmeta.updateState()
	return fmeta
}

//---------------------private helpers------------------------

func (fmeta *FileMeta) isNew() bool {
	return fmeta.mode == fileMetaModeNew
}

//updateState fetches database record and reset struct field's to contain correct value
func (fmeta *FileMeta) updateState() error {
	query := fileMetaC().Find(bson.M{"name": fmeta.name})
	count, err := query.Count()
	if err != nil {
		fmeta.err = err
		return err
	}
	if count == 0 {
		fmeta.mode = fileMetaModeNew
	} else {
		fmeta.mode = fileMetaModeNormal
		if err = query.One(&fmeta.doc); err != nil {
			fmeta.err = err
		}
	}
	return err
}

//--------------------public getter functions------------------------

//Err returns the error happened in previous operations
func (fmeta *FileMeta) Err() error {
	return fmeta.err
}

//Name returns the name
func (fmeta *FileMeta) Name() (string, error) {
	if err := fmeta.Err(); err != nil {
		return "", err
	}

	return fmeta.name, nil
}

//Size returns the file size
func (fmeta *FileMeta) Size() (uint64, error) {
	if err := fmeta.Err(); err != nil {
		return 0, err
	}
	if fmeta.isNew() {
		return 0, ErrUnsetField
	}

	return fmeta.doc.Size, nil
}

//Hash returns the file hash
func (fmeta *FileMeta) Hash() ([]byte, error) {
	if err := fmeta.Err(); err != nil {
		return nil, err
	}
	if fmeta.isNew() {
		return nil, ErrUnsetField
	}

	return fmeta.doc.Hash, nil
}

//FragCount returns the file's fragment number
func (fmeta *FileMeta) FragCount() (uint64, error) {
	if err := fmeta.Err(); err != nil {
		return 0, err
	}
	if fmeta.isNew() {
		return 0, ErrUnsetField
	}

	return fmeta.doc.FragCount, nil
}

//UploadTime returns the file's upload time
func (fmeta *FileMeta) UploadTime() (time.Time, error) {
	if err := fmeta.Err(); err != nil {
		return time.Unix(0, 0), err
	}
	if fmeta.isNew() {
		return time.Unix(0, 0), ErrUnsetField
	}

	return fmeta.doc.UploadTime, nil
}

//MerkleTree returns the file's merkle tree
func (fmeta *FileMeta) MerkleTree() (*merkle.MTree, error) {
	if err := fmeta.Err(); err != nil {
		return nil, err
	}
	if fmeta.isNew() {
		return nil, ErrUnsetField
	}

	mtree, err := merkle.Unmarshal(fmeta.doc.MerkleTree)
	if err != nil {
		return nil, err
	}

	return mtree, nil
}

//Remove removes the corresponding database record
func (fmeta *FileMeta) Remove() error {
	if err := fmeta.Err(); err != nil {
		return err
	}

	if fmeta.isNew() {
		return nil
	}
	return fileMetaC().RemoveId(fmeta.doc.ID)
}

//-----------------------editor-------------------------

//FileMetaEditor is the editing struct for edit a file meta record in db
type FileMetaEditor struct {
	fmeta *FileMeta

	doc mgoFileMeta
	err error
}

//Edit returns the editor struct and start editing
func (fmeta *FileMeta) Edit() *FileMetaEditor {
	return &FileMetaEditor{fmeta: fmeta}
}

//Err returns the happened error
func (fmeditor *FileMetaEditor) Err() error {
	return fmeditor.err
}

//SetSize sets file size
func (fmeditor *FileMetaEditor) SetSize(size uint64) *FileMetaEditor {
	if fmeditor.Err() != nil {
		return fmeditor
	}

	fmeditor.doc.Size = size
	return fmeditor
}

//SetHash sets the file hash
func (fmeditor *FileMetaEditor) SetHash(hash []byte) *FileMetaEditor {
	if fmeditor.Err() != nil {
		return fmeditor
	}

	fmeditor.doc.Hash = hash
	return fmeditor
}

//SetFragCount sets the file fragments total count
func (fmeditor *FileMetaEditor) SetFragCount(fragCount uint64) *FileMetaEditor {
	if fmeditor.Err() != nil {
		return fmeditor
	}

	fmeditor.doc.FragCount = fragCount
	return fmeditor
}

//SetUploadTime sets the uploading time of the file
func (fmeditor *FileMetaEditor) SetUploadTime(t time.Time) *FileMetaEditor {
	if fmeditor.Err() != nil {
		return fmeditor
	}

	fmeditor.doc.UploadTime = t
	return fmeditor
}

//SetMerkleTree sets the file merkle tree
func (fmeditor *FileMetaEditor) SetMerkleTree(mtree *merkle.MTree) *FileMetaEditor {
	if fmeditor.Err() != nil {
		return fmeditor
	}

	b, err := merkle.Marshal(mtree)
	if err != nil {
		fmeditor.err = err
		return fmeditor
	}
	fmeditor.doc.MerkleTree = b
	return fmeditor
}

//Done commits the changes to database
func (fmeditor *FileMetaEditor) Done() error {
	fmeditor.doc.Name = fmeditor.fmeta.name
	if err := fmeditor.Err(); err != nil {
		return err
	}
	_, err := fileMetaC().Upsert(bson.M{"name": fmeditor.fmeta.name}, &fmeditor.doc)
	fmeditor.fmeta.updateState()
	return err
}
