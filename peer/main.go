package main

import (
	"github.com/Zumium/fyer/cfg"
	db_peer "github.com/Zumium/fyer/db/peer"
	"github.com/Zumium/fyer/fragmngr"
	"github.com/Zumium/fyer/peeregister"
	rpc_peer "github.com/Zumium/fyer/rpc/peer"
	util_bin "github.com/Zumium/fyer/util/bin"
)

func main() {
	if err := cfg.Init(); err != nil {
		util_bin.ReportErr(err)
		return
	}
	if err := db_peer.Open(); err != nil {
		util_bin.ReportErr(err)
		return
	}
	defer db_peer.Close()
	if err := fragmngr.Init(); err != nil {
		util_bin.ReportErr(err)
		return
	}
	if err := peeregister.RegisterLocal(); err != nil {
		util_bin.ReportErr(err)
		return
	}
	if err := rpc_peer.Start(); err != nil {
		util_bin.ReportErr(err)
		return
	}
}
