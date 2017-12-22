package peer

import (
	pb_peer "github.com/Zumium/fyer/protos/peer"
)

//Frag represent detail info about a frag
type Frag struct {
	Index uint64 `json:"index"`
	Start int64  `json:"start"`
	Size  int64  `json:"size"`
}

func FragCommonToPb(frag Frag) *pb_peer.Frag {
	return &pb_peer.Frag{Index: frag.Index, Start: frag.Start, Size: frag.Size}
}

func FragPbToCommon(frag *pb_peer.Frag) Frag {
	return Frag{Index: frag.Index, Start: frag.Start, Size: frag.Size}
}
