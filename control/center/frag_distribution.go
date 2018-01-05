package center

import (
	db_center "github.com/Zumium/fyer/db/center"
	pb_center "github.com/Zumium/fyer/protos/center"
	"golang.org/x/net/context"
	"fmt"
)

type FragDistributionController struct{}

func (c *FragDistributionController) FragDistribution(ctx context.Context, in *pb_center.FragDistributionRequest) (*pb_center.FragDistributionResponse, error) {
	fmt.Printf("new frag distribution request: %s\n", in.String())

	handler, err := db_center.ToFragFile(in.GetName())
	if err != nil {
		return nil, err
	}
	peerList := handler.PeerList()
	resp := &pb_center.FragDistributionResponse{Distribution: make([]*pb_center.FragDistributionResponse_PeerList, 0, len(peerList))}
	for i := 0; i < len(peerList); i++ {
		resp.Distribution = append(resp.Distribution, &pb_center.FragDistributionResponse_PeerList{peerList[i]})
	}
	return resp, nil
}
