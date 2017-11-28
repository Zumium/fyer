package center

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const mgoFileMetaCollection = "filemeta"

type mgoFileMeta struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Name       string        `bson:"name,omitempty"`
	Size       uint64        `bson:"size,omitempty"`
	Hash       []byte        `bson:"hash,omitempty"`
	FragCount  uint64        `bson:"frag_count,omitempty"`
	UploadTime time.Time     `bson:"upload_time,omitempty"`
	MerkleTree []byte        `bson:"merkle_tree,omitempty"`
}

const mgoFragCollection = "frag"

type mgoFrag struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string        `bson:"name,omitempty"` //index
	FragIndex uint64        `bson:"frag_index,omitempty"`
	Peer      []string      `bson:"peer,omitempty"`
}

const mgoPeerCollection = "peer"

type mgoPeer struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	PeerID  string        `bson:"peer_id,omitempty"`
	Address string        `bson:"address,omitempty"`
}
