package center

import (
	"sync"
	"time"

	db_center "github.com/Zumium/fyer/db/center"
	"github.com/Zumium/fyer/merkle"
	pb_center "github.com/Zumium/fyer/protos/center"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

//FileRegisterStoreFileInfo is used to store file info into database
//it represents an abstract db operating process
// type FileRegisterStoreFileInfo interface {
// 	Store(name string, size uint64, hash []byte, fragCount uint64, mtree *merkle.MTree)
// }

var (
	fileRegisterOnce      sync.Once
	fileRegisterSingleton *FileRegister
)

//FileRegister -- file registering process controller
type FileRegister struct{}

//FileRegisterInstance returns the singleton instance of FileRegister
func FileRegisterInstance() *FileRegister {
	fileRegisterOnce.Do(func() {
		fileRegisterSingleton = &FileRegister{}
	})
	return fileRegisterSingleton
}

//GRPCHandler returns a GRPC server implementation function, which can be used an
//callback function to input file registering request into the controller
func (fr *FileRegister) GRPCHandler() func(context.Context, *pb_center.RegisterRequest) (*google_protobuf.Empty, error) {
	return func(ctx context.Context, in *pb_center.RegisterRequest) (*google_protobuf.Empty, error) {
		dbHandler, err := db_center.ToFileMeta(in.Name)
		if err != nil {
			return nil, err
		}
		mtree, err := merkle.Unmarshal(in.MerkleTree)
		if err != nil {
			return nil, err
		}
		editor := dbHandler.Edit()
		editor.SetSize(in.Size).SetHash(in.Hash).SetFragCount(in.FragCount).SetMerkleTree(mtree).SetUploadTime(time.Now())
		if err := editor.Done(); err != nil {
			return nil, err
		}

		return new(google_protobuf.Empty), nil
	}
}
