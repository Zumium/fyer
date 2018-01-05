package main

import "C"
import (
	control_fyerwork "github.com/Zumium/fyer/control/fyerwork"
	"github.com/Zumium/fyer/filemanager"
	rpc_fyerwork "github.com/Zumium/fyer/rpc/fyerwork"
)

var localServeAddr string

//------------------------------ Configuration --------------------------

//export SetLocalServeAddress
func SetLocalServeAddress(laddr string) {
	localServeAddr = laddr
}

//-----------------------------------------------------------------------

//------------------------------ File Management ------------------------

//export RegisterFile
func RegisterFile(fpath string) int {
	if err := filemanager.Register(fpath); err != nil {
		return -1
	}
	return 0
}

//export UnregisterFile
func UnregisterFile(name string) bool {
	return filemanager.Unregister(name)
}

//-----------------------------------------------------------------------

//----------------------------- Fyerwork RPC Server ---------------------

//export StartFyerworkServer
func StartFyerworkServer() int {
	if err := rpc_fyerwork.Start(); err != nil {
		return -1
	}
	return 0
}

//export StartFyerworkServerInBackground
func StartFyerworkServerInBackground() {
	rpc_fyerwork.StartInBackground()
}

//export WaitFyerworkServer
func WaitFyerworkServer() {
	rpc_fyerwork.Wait()
}

//-----------------------------------------------------------------------

//----------------------------- Upload File -----------------------------

//export UploadFile
func UploadFile(name string, size uint64, hash []byte) int {
	if err := control_fyerwork.UploadFile(name, size, hash, localServeAddr); err != nil {
		return -1
	}
	return 0
}

//-----------------------------------------------------------------------

//----------------------------- Download File ---------------------------

//export DownloadFile
func DownloadFile(name string, storePath string) int {
	if err := control_fyerwork.Download(name, storePath); err != nil {
		return -1
	}
	return 0
}

func main() {}
