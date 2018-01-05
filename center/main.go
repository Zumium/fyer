package main

import (
	"github.com/Zumium/fyer/cfg"
	db_center "github.com/Zumium/fyer/db/center"
	"github.com/Zumium/fyer/peeregister"
	rpc_center "github.com/Zumium/fyer/rpc/center"
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
	if err := db_center.Init(); err != nil {
		util_bin.ReportErr(err)
		return
	}
	defer db_center.Close()

	fmt.Println("Initializing peer registration server")
	if err := peeregister.InitServer(); err != nil {
		util_bin.ReportErr(err)
		return
	}

	fmt.Println("Starting center RPC server")
	if err := rpc_center.Start(); err != nil {
		util_bin.ReportErr(err)
		return
	}
}
