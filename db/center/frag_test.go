package center

import (
	"github.com/Zumium/fyer/common"
	"github.com/spf13/viper"
	"testing"
)

func TestFragOperations(t *testing.T) {
	viper.Set("mongo_address", "127.0.0.1:27017")

	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Close()

	frags := []common.Frag{
		common.Frag{Index: 0, Start: 0, Size: 3},
		common.Frag{Index: 1, Start: 3, Size: 3},
		common.Frag{Index: 2, Start: 6, Size: 2},
	}

	peerList := [][]string{
		[]string{"A", "B", "C"},
		[]string{"B", "A", "C"},
		[]string{"C", "B", "A"},
	}

	handler, err := ToFragFile("hello")
	if err != nil {
		t.Fatal(err)
	}
	if err := handler.Edit().SetFrags(frags).SetPeerList(peerList).Done(); err != nil {
		t.Fatal(err)
	}
	defer handler.Remove()

	handler2, err := ToFragFile("hello")
	if err != nil {
		t.Fatal(err)
	}

	frags2 := handler2.Frags()
	if f := frags2[0]; f.Index != 0 || f.Start != 0 || f.Size != 3 {
		t.Fatal("frag 0 data haven't been correctly stored")
	}
	if f := frags2[1]; f.Index != 1 || f.Start != 3 || f.Size != 3 {
		t.Fatal("frag 1 data haven't been correctly stored")
	}
	if f := frags2[2]; f.Index != 2 || f.Start != 6 || f.Size != 2 {
		t.Fatal("frag 2 data haven't been correctly stored")
	}

	peerList2 := handler2.PeerList()
	if l := peerList2[0]; l[0] != "A" || l[1] != "B" || l[2] != "C" {
		t.Fatal("peer list 0 haven't been correctly stored")
	}
	if l := peerList2[1]; l[0] != "B" || l[1] != "A" || l[2] != "C" {
		t.Fatal("peer list 1 haven't been correctly stored")
	}
	if l := peerList2[2]; l[0] != "C" || l[1] != "B" || l[2] != "A" {
		t.Fatal("peer list 2 haven't been correctly stored")
	}
}
