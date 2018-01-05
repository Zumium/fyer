package main

import "C"
import (
	control_fyerwork "github.com/Zumium/fyer/control/fyerwork"
	"github.com/Zumium/fyer/filemanager"
	rpc_fyerwork "github.com/Zumium/fyer/rpc/fyerwork"
)

var localServeAddr string
var messageSize = 128 * 1024 * 1024

//------------------------------ Configuration --------------------------

//export SetLocalServeAddress
func SetLocalServeAddress(laddr string) {
	localServeAddr = laddr
}

//export SetMessageSize
func SetMessageSize(size int) {
	messageSize = size
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
func StartFyerworkServer(port int) int {
	if err := rpc_fyerwork.Start(port, messageSize); err != nil {
		return -1
	}
	return 0
}

//export StartFyerworkServerInBackground
func StartFyerworkServerInBackground(port int) {
	rpc_fyerwork.StartInBackground(port, messageSize)
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
