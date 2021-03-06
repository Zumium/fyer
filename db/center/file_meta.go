package center

import (
	"time"

	//"github.com/Zumium/fyer/merkle"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//fileMetaC returns the file meta collection
func fileMetaC() *mgo.Collection {
	return fileMetaCollection
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

//Files returns all files' names from db
func Files() ([]string, error) {
	var results []struct {
		Name string `bson:"name"`
	}
	if err := fileMetaC().Find(nil).Select(bson.M{"name": 1, "_id": 0}).All(&results); err != nil {
		return nil, err
	}

	names := make([]string, 0, len(results))
	for _, r := range results {
		names = append(names, r.Name)
	}

	return names, nil
}

//ToFileMeta creates a new FileMeta to apply furthur db operaions
//It garantees that the error is clear if instance has been successfully created
func ToFileMeta(name string) (*FileMeta, error) {
	fmeta := &FileMeta{name: name}
	return fmeta, fmeta.updateState()
}

//---------------------public helpers------------------------

//IsNew returns whether the record exists in database already
func (fmeta *FileMeta) IsNew() bool {
	return fmeta.mode == fileMetaModeNew
}

//Err returns the latest error occured
func (fmeta *FileMeta) Err() error {
	return fmeta.err
}

//ClearErr resets the internal error to nil
func (fmeta *FileMeta) ClearErr() {
	fmeta.err = nil
}

//---------------------private helpers------------------------

//updateState fetches database record and reset struct field's to contain correct value
func (fmeta *FileMeta) updateState() error {
	query := fileMetaC().Find(bson.M{"name": fmeta.name})
	count, err := query.Count()
	if err != nil {
		return err
	}

	if count == 0 {
		fmeta.mode = fileMetaModeNew
	} else {
		fmeta.mode = fileMetaModeNormal
		err = query.One(&fmeta.doc)
	}
	return err
}

//--------------------public getter functions------------------------

//Name returns the name
func (fmeta *FileMeta) Name() string {
	return fmeta.name
}

//Size returns the file size
func (fmeta *FileMeta) Size() uint64 {
	return fmeta.doc.Size
}

//Hash returns the file hash
func (fmeta *FileMeta) Hash() []byte {
	return fmeta.doc.Hash
}

//FragCount returns the file's fragment number
func (fmeta *FileMeta) FragCount() uint64 {
	return fmeta.doc.FragCount
}

//UploadTime returns the file's upload time
func (fmeta *FileMeta) UploadTime() time.Time {
	return fmeta.doc.UploadTime
}

//RawMerkleTree returns the marshaled merkle tree
//func (fmeta *FileMeta) RawMerkleTree() []byte {
//	return fmeta.doc.MerkleTree
//}
//
////MerkleTree returns the file's merkle tree
//func (fmeta *FileMeta) MerkleTree() *merkle.MTree {
//	mtree, err := merkle.Unmarshal(fmeta.doc.MerkleTree)
//	if err != nil {
//		fmeta.err = err
//		return nil
//	}
//
//	return mtree
//}

//Remove removes the corresponding database record
func (fmeta *FileMeta) Remove() error {
	return fileMetaC().RemoveId(fmeta.doc.ID)
}

//-----------------------editor-------------------------

//FileMetaEditor is the editing struct for edit a file meta record in db
type FileMetaEditor struct {
	fmeta *FileMeta

	doc bson.M
	err error
}

//Edit returns the editor struct and start editing
func (fmeta *FileMeta) Edit() *FileMetaEditor {
	return &FileMetaEditor{fmeta: fmeta, doc: bson.M{}}
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

	//fmeditor.doc.Size = size
	fmeditor.doc["size"] = size
	return fmeditor
}

//SetHash sets the file hash
func (fmeditor *FileMetaEditor) SetHash(hash []byte) *FileMetaEditor {
	if fmeditor.Err() != nil {
		return fmeditor
	}

	//fmeditor.doc.Hash = hash
	fmeditor.doc["hash"] = hash
	return fmeditor
}

//SetFragCount sets the file fragments total count
func (fmeditor *FileMetaEditor) SetFragCount(fragCount uint64) *FileMetaEditor {
	if fmeditor.Err() != nil {
		return fmeditor
	}

	//fmeditor.doc.FragCount = fragCount
	fmeditor.doc["frag_count"] = fragCount
	return fmeditor
}

//SetUploadTime sets the uploading time of the file
func (fmeditor *FileMetaEditor) SetUploadTime(t time.Time) *FileMetaEditor {
	if fmeditor.Err() != nil {
		return fmeditor
	}

	//fmeditor.doc.UploadTime = t
	fmeditor.doc["upload_time"] = t
	return fmeditor
}

//Done commits the changes to database
func (fmeditor *FileMetaEditor) Done() error {
	if err := fmeditor.Err(); err != nil {
		return err
	}
	//fmeditor.doc.Name = fmeditor.fmeta.name
	//if _, err := fileMetaC().Upsert(bson.M{"name": fmeditor.fmeta.name}, &fmeditor.doc); err != nil {
	if _, err := fileMetaC().Upsert(&mgoFileMeta{Name: fmeditor.fmeta.name}, bson.M{"$set": fmeditor.doc}); err != nil {
		return err
	}
	return fmeditor.fmeta.updateState()
}
