package peer

import (
	"github.com/Zumium/fyer/cfg"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	//dbInstance is to hold the global dababase instance
	dbInstance *leveldb.DB
)

//Open initializes the database which is located by configs
func Open() (err error) {
	dbInstance, err = leveldb.OpenFile(cfg.DBPath(), nil)
	return
}

//Close closes the database
func Close() error {
	return instance().Close()
}

//instance (for internal use) returns the global database instance
func instance() *leveldb.DB {
	return dbInstance
}
