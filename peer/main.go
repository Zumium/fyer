package main

import (
	"github.com/Zumium/fyer/cfg"
	db_peer "github.com/Zumium/fyer/db/peer"
	"github.com/Zumium/fyer/fragmngr"
	"github.com/Zumium/fyer/peeregister"
	rpc_peer "github.com/Zumium/fyer/rpc/peer"
	util_bin "github.com/Zumium/fyer/util/bin"
	"fmt"
)

func main() {
	fmt.Println("Initializing config")
	if err := cfg.Init(); err != nil {
		util_bin.ReportErr(err)
		return
	}
	//if err := log.Init(); err != nil {
	//	util_bin.ReportErr(err)
	//	return
	//}
	fmt.Println("Initializing center database")
	if err := db_peer.Open(); err != nil {
		util_bin.ReportErr(err)
		return
	}
	defer db_peer.Close()

	fmt.Println("Initializing fragment manager")
	if err := fragmngr.Init(); err != nil {
		util_bin.ReportErr(err)
		return
	}

	fmt.Println("Registering self to center")
	if err := peeregister.RegisterLocal(); err != nil {
		util_bin.ReportErr(err)
		return
	}

	fmt.Println("Starting peer RPC server")
	if err := rpc_peer.Start(); err != nil {
		util_bin.ReportErr(err)
		return
	}
}
