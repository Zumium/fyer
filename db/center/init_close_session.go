package center

import (
	"github.com/Zumium/fyer/cfg"
	"gopkg.in/mgo.v2"
)

var (
	session *mgo.Session
)

func mgoSession() *mgo.Session {
	return session
}

func mgoFyerDB() *mgo.Database {
	return mgoSession().DB("fyer")
}

func ensureIndexes() (err error) {
	if err = session.DB("fyer").C(mgoFileMetaCollection).EnsureIndexKey("name"); err != nil {
		return
	}
	if err = session.DB("fyer").C(mgoFragCollection).EnsureIndexKey("name"); err != nil {
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
	if err = ensureIndexes(); err != nil {
		return
	}
	return
}

//Close closes the mongo session
func Close() (err error) {
	session.Close()
	return
}
