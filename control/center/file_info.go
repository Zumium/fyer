package center

import (
	"fmt"
	db_center "github.com/Zumium/fyer/db/center"
	pb_center "github.com/Zumium/fyer/protos/center"
	"golang.org/x/net/context"
)

//FileInfoController handles the FileInfoController RPC process
type FileInfoController struct{}

//handles GRPC request
func (finfo *FileInfoController) FileInfo(ctx context.Context, in *pb_center.FileInfoRequest) (*pb_center.FileInfoResponse, error) {
	fmt.Printf("new file info request: %s\n", in.String())

	dbHandler, err := db_center.ToFileMeta(in.GetName())
	if err != nil {
		return nil, err
	}
	finfoResp := &pb_center.FileInfoResponse{
		Size:      dbHandler.Size(),
		Hash:      dbHandler.Hash(),
		FragCount: dbHandler.FragCount(),
		//MerkleTree: dbHandler.RawMerkleTree(),
	}
	if err := dbHandler.Err(); err != nil {
		return nil, err
	}
	return finfoResp, nil
}
