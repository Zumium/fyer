package center

import (
	"fmt"
	db_center "github.com/Zumium/fyer/db/center"
	pb_center "github.com/Zumium/fyer/protos/center"
	google_protobuf_empty "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

type FilesController struct{}

func (controller *FilesController) Files(ctx context.Context, in *google_protobuf_empty.Empty) (*pb_center.FilesResponse, error) {
	fmt.Println("new files request accepted")

	names, err := db_center.Files()
	if err != nil {
		return nil, err
	}
	return &pb_center.FilesResponse{Files: names}, nil
}
