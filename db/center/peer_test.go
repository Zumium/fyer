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
