package center

import (
	"sync"

	db_center "github.com/Zumium/fyer/db/center"
	pb_center "github.com/Zumium/fyer/protos/center"
	"golang.org/x/net/context"
)

var (
	fileInfoOnce      sync.Once
	fileInfoSingleton *FileInfo
)

//FileInfo handles the FileInfo RPC process
type FileInfo struct{}

//FileInfoInstance returns the singleton instance of FileInfo
func FileInfoInstance() *FileInfo {
	fileInfoOnce.Do(func() {
		fileInfoSingleton = &FileInfo{}
	})
	return fileInfoSingleton
}

//GRPCHandler returns a grpc handler function
func (finfo *FileInfo) GRPCHandler() func(context.Context, *pb_center.FileInfoRequest) (*pb_center.FileInfoResponse, error) {
	return func(ctx context.Context, in *pb_center.FileInfoRequest) (*pb_center.FileInfoResponse, error) {
		dbHandler, err := db_center.ToFileMeta(in.GetName())
		if err != nil {
			return nil, err
		}
		finfoResp := &pb_center.FileInfoResponse{
			Size:       dbHandler.Size(),
			Hash:       dbHandler.Hash(),
			FragCount:  dbHandler.FragCount(),
			MerkleTree: dbHandler.RawMerkleTree(),
		}
		if err := dbHandler.Err(); err != nil {
			return nil, err
		}
		return finfoResp, nil
	}
}
