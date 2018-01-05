package center

import (
	"encoding/json"
	"fmt"
	db_center "github.com/Zumium/fyer/db/center"
	pb_center "github.com/Zumium/fyer/protos/center"
	"golang.org/x/net/context"
)

type FragInfoController struct{}

func (c *FragInfoController) FragInfo(ctx context.Context, in *pb_center.FragInfoRequest) (*pb_center.FragInfoResponse, error) {
	fmt.Printf("new frag info request: %s\n", in.String())

	handler, err := db_center.ToFragFile(in.GetName())
	if err != nil {
		return nil, err
	}
	bFrags, err := json.Marshal(handler.Frags())
	if err != nil {
		return nil, err
	}
	return &pb_center.FragInfoResponse{bFrags}, nil
}
