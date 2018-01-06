package center

import (
	"github.com/Zumium/fyer/cfg"
	"gopkg.in/mgo.v2"
)

var (
	session *mgo.Session
	db *mgo.Database
	fileMetaCollection, fragCollection, peerCollection *mgo.Collection
)

//------------------------------------------------------------------------------

func doCreates() {
	db = session.DB("fyer")
	fileMetaCollection = db.C(mgoFileMetaCollection)
	fragCollection = db.C(mgoFragCollection)
	peerCollection = db.C(mgoPeerCollection)
}

func ensureIndexes() (err error) {
	if err = fileMetaCollection.EnsureIndexKey("name"); err != nil {
		return
	}
	if err = fragCollection.EnsureIndexKey("name"); err != nil {
		return
	}
	if err = peerCollection.EnsureIndexKey("peer_id");err!=nil{
		return
	}
	return
}

//Init initializes mongodb session
func Init() (err error) {
	session, err = mgo.Dial(cfg.MongoAddress())
	if err != nil {
		return
	}
	doCreates()
	if err = ensureIndexes(); err != nil {
		return
	}
	return
}

//Close closes the mongo session
func Close() error {
	session.Close()
	return nil
}
