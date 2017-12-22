package center

import (
	db_center "github.com/Zumium/fyer/db/center"
)

func ResolvePeerID(peerID string) (string, error) {
	peer, err := db_center.ToPeerID(peerID)
	if err != nil {
		return peerID, err
	}
	return peer.Address(), nil
}
