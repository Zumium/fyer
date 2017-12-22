package peer

import (
	util_bin "github.com/Zumium/fyer/util/bin"
	"github.com/Zumium/fyer/cfg"
	db_peer "github.com/Zumium/fyer/db/peer"
	"github.com/Zumium/fyer/fragmngr"
	rpc_peer "github.com/Zumium/fyer/rpc/peer"
)

func main() {
	if err:=cfg.Init();err!=nil {
		util_bin.ReportErr(err)
		return
	}
	if err:=db_peer.Open(); err!=nil {
		util_bin.ReportErr(err)
		return
	}
	defer db_peer.Close()
	if err:=fragmngr.Init(); err!=nil {
		util_bin.ReportErr(err)
		return
	}
	if err:=rpc_peer.Start();err!=nil {
		util_bin.ReportErr(err)
		return
	}
}
