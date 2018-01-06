package main

import "C"
import (
	control_fyerwork "github.com/Zumium/fyer/control/fyerwork"
	"github.com/Zumium/fyer/filemanager"
	rpc_fyerwork "github.com/Zumium/fyer/rpc/fyerwork"
	"github.com/spf13/viper"
)

//------------------------------ Configuration --------------------------

//REQUIRED
//export SetLocalServeAddress
func SetLocalServeAddress(laddr string) {
	viper.Set("local_serve_address", laddr)
}

//REQUIRED
//export SetCenterAddress
func SetCenterAddress(addr string) {
	viper.Set("center_adddress", addr)
}

//export SetPort
func SetPort(port int) {
	viper.Set("port", port)
}

//export SetMaxSendRecvMessageSize
func SetMaxSendRecvMessageSize(size int) {
	viper.Set("max_send_recv_msg_size", size)
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
	if err := control_fyerwork.UploadFile(name, size, hash, viper.GetString("local_serve_address")); err != nil {
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

//-----------------------------------------------------------------------

//---------------------------- Initializing -----------------------------

func init() {
	viper.Set("port", 4201)
	viper.Set("max_send_recv_msg_size", 128*1024*1024)
}

//-----------------------------------------------------------------------
