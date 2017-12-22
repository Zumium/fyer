package center

import (
	util_bin "github.com/Zumium/fyer/util/bin"
	"github.com/Zumium/fyer/cfg"
	db_center "github.com/Zumium/fyer/db/center"
	"github.com/Zumium/fyer/peeregister"
	rpc_center "github.com/Zumium/fyer/rpc/center"
)

func main() {
	if err:=cfg.Init();err!=nil {
		util_bin.ReportErr(err)
		return
	}
	if err:=db_center.Init();err!=nil {
		util_bin.ReportErr(err)
		return
	}
	defer db_center.Close()
	if err:=peeregister.InitServer(); err!=nil {
		util_bin.ReportErr(err)
		return
	}
	if err:=rpc_center.Start(); err!=nil {
		util_bin.ReportErr(err)
		return
	}
}