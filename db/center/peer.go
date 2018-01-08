package center

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func peerC() *mgo.Collection {
	return peerCollection
}

type peerRecordMode uint8

const (
	peerRecordNew peerRecordMode = iota
	peerRecordNormal
)

//Peer represents a peer database record
type Peer struct {
	peerID string

	mode peerRecordMode
	doc  mgoPeer
	err  error
}

//ToPeerID creates a Peer instance to apply futhur db operations
func ToPeerID(peerID string) (*Peer, error) {
	p := &Peer{peerID: peerID}
	return p, p.updateState()
}

//AllPeers gets all peers
func AllPeers() ([]*Peer, error) {
	var mgoPeers []mgoPeer
	if err := peerC().Find(nil).All(&mgoPeers); err != nil {
		return nil, err
	}

	peers := make([]*Peer, len(mgoPeers))
	for i, p := range mgoPeers {
		peers[i] = &Peer{p.PeerID, peerRecordNormal, p, nil}
	}

	return peers, nil
}

//---------------------public helpers------------------------

//IsNew returns whether the record exists in database already
func (p *Peer) IsNew() bool {
	return p.mode == peerRecordNew
}

//Err returns the latest error occured
func (p *Peer) Err() error {
	return p.err
}

//ClearErr resets the internal error to nil
func (p *Peer) ClearErr() {
	p.err = nil
}

//---------------------private helpers------------------------

//updateState fetches database record and reset struct field's to contain correct value
func (p *Peer) updateState() error {
	query := peerC().Find(bson.M{"peer_id": p.peerID})
	count, err := query.Count()
	if err != nil {
		return err
	}

	if count == 0 {
		p.mode = peerRecordNew
	} else {
		p.mode = peerRecordNormal
		err = query.One(&p.doc)
	}
	return err
}

//--------------------public getter functions------------------------

//PeerID returns the peer id
func (p *Peer) PeerID() string {
	return p.peerID
}

//Address returns the peer address that registered before
func (p *Peer) Address() string {
	return p.doc.Address
}

//Remove removes the corresponding database record
func (p *Peer) Remove() error {
	return peerC().RemoveId(p.doc.ID)
}

//-----------------------editor-------------------------

//PeerEditor is the editing struct for edit a peer record in db
type PeerEditor struct {
	peer *Peer

	doc bson.M
	err error
}

//Edit returns the editing structure and start editing
func (p *Peer) Edit() *PeerEditor {
	return &PeerEditor{peer: p, doc: bson.M{}}
}

//Err returns the latest happened error
func (pe *PeerEditor) Err() error {
	return pe.err
}

//SetAddress sets the address
func (pe *PeerEditor) SetAddress(address string) *PeerEditor {
	if pe.Err() != nil {
		return pe
	}

	pe.doc["address"] = address
	return pe
}

//Done commits the changes
func (pe *PeerEditor) Done() error {
	if err := pe.Err(); err != nil {
		return err
	}
	if _, err := peerC().Upsert(&mgoPeer{PeerID: pe.peer.peerID}, bson.M{"$set": pe.doc}); err != nil {
		return err
	}
	return pe.peer.updateState()
}
