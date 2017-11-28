package center

import (
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
	return mgoFyerDB().C(mgoFragCollection)
}

//Frag is file fragment database handler
type Frag struct {
	name      string
	fragIndex uint64

	mode fragRecordMode
	doc  mgoFrag
	err  error
}

//--------------------constructor-----------------------

//ToFragFile creates a new handler of given file
func ToFragFile(name string) *Frag {
	return &Frag{name: name}
}

//SelectIndex selects the given frag index of file
func (f *Frag) SelectIndex(idx uint64) error {
	f.fragIndex = idx
	return f.updateState()
}

//--------------------private helpers--------------------

func (f *Frag) updateState() error {
	query := fragC().Find(bson.M{"name": f.name, "index": f.fragIndex})
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

//FragIndex returns the index of the frag
func (f *Frag) FragIndex() uint64 {
	return f.doc.FragIndex
}

//Peer returns the list of peers in which stores the specified frag
func (f *Frag) Peer() []string {
	return f.doc.Peer
}

//====================editor=======================

//FragEditor the editing handler for frag record
type FragEditor struct {
	frag *Frag

	doc mgoFrag
	err error
}

//Edit returns the editor instance and start editing
func (f *Frag) Edit() *FragEditor {
	return &FragEditor{frag: f}
}

//Err returns the latest error that occured
func (fe *FragEditor) Err() error {
	return fe.err
}

//SetPeer sets the list of peers who stores the given frag
func (fe *FragEditor) SetPeer(peer []string) *FragEditor {
	if err := fe.Err(); err != nil {
		return fe
	}

	fe.doc.Peer = peer
	return fe
}

//Done commits the changes to database
func (fe *FragEditor) Done() error {
	if err := fe.Err(); err != nil {
		return err
	}
	fe.doc.Name, fe.doc.FragIndex = fe.frag.name, fe.frag.fragIndex
	if _, err := fragC().Upsert(bson.M{"name": fe.frag.name, "frag_index": fe.frag.fragIndex}, fe.doc); err != nil {
		return err
	}
	return fe.frag.updateState()
}
