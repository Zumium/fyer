package main

/*
#include <stdlib.h>
#include <stdint.h>
#include <string.h>

typedef char* char_ptr;

typedef struct {
	uint64_t size;
	unsigned char hash[16];
	uint64_t frag_count;
	int64_t upload_time;
} file_info;

void set_c_str_array(char *array[], int index, char *str);
*/
import "C"
import (
	control_fyerwork "github.com/Zumium/fyer/control/fyerwork"
	"github.com/Zumium/fyer/filemanager"
	rpc_fyerwork "github.com/Zumium/fyer/rpc/fyerwork"
	"github.com/spf13/viper"
	"unsafe"
)

//------------------------------ Configuration --------------------------

//REQUIRED
//export SetLocalServeAddress
func SetLocalServeAddress(laddr *C.char) {
	viper.Set("local_serve_address", C.GoString(laddr))
}

//REQUIRED
//export SetCenterAddress
func SetCenterAddress(addr *C.char) {
	viper.Set("center_adddress", C.GoString(addr))
}

//export SetPort
func SetPort(port C.int) {
	viper.Set("port", int(port))
}

//export SetMaxSendRecvMessageSize
func SetMaxSendRecvMessageSize(size C.int) {
	viper.Set("max_send_recv_msg_size", int(size))
}

//-----------------------------------------------------------------------

//------------------------------ File Management ------------------------

//export RegisterFile
func RegisterFile(fpath *C.char) C.int {
	if err := filemanager.Register(C.GoString(fpath)); err != nil {
		return -1
	}
	return 0
}

//export UnregisterFile
func UnregisterFile(name *C.char) C.int {
	if filemanager.Unregister(C.GoString(name)) {
		return 1
	}
	return 0
}

//-----------------------------------------------------------------------

//----------------------------- Fyerwork RPC Server ---------------------

//export StartFyerworkServer
func StartFyerworkServer() C.int {
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

//----------------------------- Querying Operations ---------------------

//export Files
func Files(buf ***C.char, bufLen *C.int) C.int {
	names, err := control_fyerwork.Files()
	if err != nil {
		return -1
	}

	l := len(names)
	*buf = (**C.char)(C.calloc(C.size_t(l), C.sizeof_char_ptr))
	*bufLen = C.int(l)
	for i, name := range names {
		C.set_c_str_array(*buf, C.int(i), C.CString(name))
	}

	return 0
}

//export FileInfo
func FileInfo(name *C.char, finfo *C.file_info) C.int {
	size,hash,fragCount,uploadTime,err:=control_fyerwork.FileInfo(C.GoString(name))
	if err != nil {
		return -1
	}

	finfo.size = C.uint64_t(size)

	cHash := C.CBytes(hash)
	defer C.free(cHash)
	C.memcpy(unsafe.Pointer(&finfo.hash[0]), cHash, 16)

	finfo.frag_count = C.uint64_t(fragCount)

	finfo.upload_time = C.int64_t(uploadTime.Unix())

	return 0;
}

//-----------------------------------------------------------------------

//----------------------------- Upload File -----------------------------

//export UploadFile
func UploadFile(name *C.char, size C.uint64_t, hash unsafe.Pointer, hashLen C.int) C.int {
	if err := control_fyerwork.UploadFile(C.GoString(name), uint64(size), C.GoBytes(hash, hashLen), viper.GetString("local_serve_address")); err != nil {
		return -1
	}
	return 0
}

//-----------------------------------------------------------------------

//----------------------------- Download File ---------------------------

//export DownloadFile
func DownloadFile(name *C.char, storePath *C.char) C.int {
	if err := control_fyerwork.Download(C.GoString(name), C.GoString(storePath)); err != nil {
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
