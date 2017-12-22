package center

import (
	"testing"

	"github.com/spf13/viper"
)

func TestPeerOperations(t *testing.T) {
	viper.Set("mongo_address", "127.0.0.1:27017")

	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Close()

	p, _ := ToPeerID("test")
	if err := p.Edit().SetAddress("1.2.3.4").Done(); err != nil {
		t.Fatal(err)
	}
	defer p.Remove()

	p2, _ := ToPeerID("test")
	addr := p2.Address()
	if addr != "1.2.3.4" {
		t.Fatalf("frag count mismatching: is %s, should be %s\n", addr, "1.2.3.4")
	}
}

func TestAllPeers(t *testing.T) {
	viper.Set("mongo_address", "127.0.0.1:27017")

	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Close()

	pA, err := ToPeerID("A")
	if err != nil {
		t.Fatal(err)
	}
	if err := pA.Edit().SetAddress("192.168.1.1").Done(); err != nil {
		t.Fatal(err)
	}
	defer pA.Remove()

	pB, err := ToPeerID("B")
	if err != nil {
		t.Fatal(err)
	}
	if err := pB.Edit().SetAddress("192.168.1.2").Done(); err != nil {
		t.Fatal(err)
	}
	defer pB.Remove()

	peers, err := AllPeers()
	if err != nil {
		t.Fatal(err)
	}
	peerID2Address := map[string]string{"A": "192.168.1.1", "B": "192.168.1.2"}
	for _, p := range peers {
		if p.Address() != peerID2Address[p.PeerID()] {
			t.Fatalf("data error: %s-%s", p.PeerID(), p.Address())
		}
	}
}
