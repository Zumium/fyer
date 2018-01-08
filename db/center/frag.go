package center

import (
	common "github.com/Zumium/fyer/common"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type fragRecordMode uint8

const (
	fragRecordModeNew fragRecordMode = iota
	fragRecordModeNormal
)

//fragC returns the collection of file frag infos
func fragC() *mgo.Collection {
	return fragCollection
}

//Frag is file fragment database handler
type Frag struct {
	name string
	//fragIndex uint64

	mode fragRecordMode
	doc  mgoFrag
	err  error
}

//--------------------constructor-----------------------

//ToFragFile creates a new handler of given file
func ToFragFile(name string) (*Frag, error) {
	frag := &Frag{name: name}
	return frag, frag.updateState()
}

//--------------------private helpers--------------------

func (f *Frag) updateState() error {
	query := fragC().Find(bson.M{"name": f.name})
	count, err := query.Count()
	if err != nil {
		return err
	}
	if count == 0 {
		f.mode = fragRecordModeNew
	} else {
		f.mode = fragRecordModeNormal
		err = query.One(&f.doc)
	}
	return err
}

//--------------------public helpers--------------------

//IsNew returns whether the record exists in database already
func (f *Frag) IsNew() bool {
	return f.mode == fragRecordModeNew
}

//Err returns the latest error occured
func (f *Frag) Err() error {
	return f.err
}

//ClearErr resets the internal error to nil
func (f *Frag) ClearErr() {
	f.err = nil
}

//-----------------public getters--------------------------

//Name returns frag file's name
func (f *Frag) Name() string {
	return f.doc.Name
}

//Frags returns the frags
func (f *Frag) Frags() []common.Frag {
	return f.doc.Frags
}

//FragByIndex returns the frag with specified index
func (f *Frag) FragByIndex(idx uint64) common.Frag {
	return f.Frags()[idx]
}

//PeerList returns the list of peers in which stores the specified frag
func (f *Frag) PeerList() [][]string {
	return f.doc.PeerList
}

func (f *Frag) Remove() error {
	return fragC().Remove(bson.M{"name": f.name})
}

//====================editor=======================

//FragEditor the editing handler for frag record
type FragEditor struct {
	frag *Frag

	doc bson.M
	err error
}

//Edit returns the editor instance and start editing
func (f *Frag) Edit() *FragEditor {
	return &FragEditor{frag: f, doc: bson.M{}}
}

//Err returns the latest error that occured
func (fe *FragEditor) Err() error {
	return fe.err
}

//SetPeerList sets the list of peers
func (fe *FragEditor) SetPeerList(peerList [][]string) *FragEditor {
	if err := fe.Err(); err != nil {
		return fe
	}

	fe.doc["peer_list"] = peerList
	return fe
}

//SetFrags sets the frags
func (fe *FragEditor) SetFrags(frags []common.Frag) *FragEditor {
	if err := fe.Err(); err != nil {
		return fe
	}

	fe.doc["frags"] = frags
	return fe
}

//Done commits the changes to database
func (fe *FragEditor) Done() error {
	if err := fe.Err(); err != nil {
		return err
	}
	if _, err := fragC().Upsert(&mgoFrag{Name: fe.frag.name}, bson.M{"$set": fe.doc}); err != nil {
		return err
	}
	return fe.frag.updateState()
}
