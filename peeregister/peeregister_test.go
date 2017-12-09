package peeregister

import (
	"github.com/Zumium/fyer/cfg"
	db_center "github.com/Zumium/fyer/db/center"
	"github.com/spf13/viper"
	"testing"
)

func TestPeerRegister(t *testing.T) {
	//setup config
	viper.Set("peer_id", "A")
	viper.Set("center_address", "127.0.0.1")
	viper.Set("peer_register_port", 4103)
	//setup database
	if err := db_center.Init(); err != nil {
		t.Fatal(err)
	}
	defer db_center.Close()
	//start server
	if err := InitServer(); err != nil {
		t.Fatal(err)
	}

	if err := RegisterLocal(); err != nil {
		t.Fatal(err)
	}

	dbHandler, err := db_center.ToPeerID(cfg.PeerID())
	if err != nil {
		t.Fatal(err)
	}
	address := dbHandler.Address()
	if err := dbHandler.Err(); err != nil {
		t.Fatal(err)
	}
	if address != "127.0.0.1" {
		t.Fatal("registered address not matching")
	}
}
