package center

import (
	db_center "github.com/Zumium/fyer/db/center"
	pb_center "github.com/Zumium/fyer/protos/center"
	"golang.org/x/net/context"
	"fmt"
)

//PeerInfoController handles the PeerInfoController RPC process
type PeerInfoController struct{}

//handles GRPC request
func (pinfo *PeerInfoController) PeerInfo(ctx context.Context, in *pb_center.PeerInfoRequest) (*pb_center.PeerInfoResponse, error) {
	fmt.Printf("new peer info request: %s\n", in.String())

	dbHandler, err := db_center.ToPeerID(in.GetPeerId())
	if err != nil {
		return nil, err
	}
	pinfoResp := &pb_center.PeerInfoResponse{}
	if dbHandler.IsNew() {
		pinfoResp.Address = in.GetPeerId()
	} else {
		pinfoResp.Address = dbHandler.Address()
		if dbHandler.Err() != nil {
			return nil, dbHandler.Err()
		}
	}
	return pinfoResp, nil
}
