package center

import (
	"sync"

	db_center "github.com/Zumium/fyer/db/center"
	pb_center "github.com/Zumium/fyer/protos/center"
	"golang.org/x/net/context"
)

var (
	peerInfoOnce      sync.Once
	peerInfoSingleton *PeerInfo
)

//PeerInfo handles the PeerInfo RPC process
type PeerInfo struct{}

//PeerInfoInstance returns the singleton instance of PeerInfo
func PeerInfoInstance() *PeerInfo {
	peerInfoOnce.Do(func() {
		peerInfoSingleton = &PeerInfo{}
	})
	return peerInfoSingleton
}

//GRPCHandler returns a grpc handler function
func (pinfo *PeerInfo) GRPCHandler() func(context.Context, *pb_center.PeerInfoRequest) (*pb_center.PeerInfoResponse, error) {
	return func(ctx context.Context, in *pb_center.PeerInfoRequest) (*pb_center.PeerInfoResponse, error) {
		dbHandler, err := db_center.ToPeerID(in.GetPeerID())
		if err != nil {
			return nil, err
		}
		pinfoResp := &pb_center.PeerInfoResponse{
			Address: dbHandler.Address(),
		}
		if err := dbHandler.Err(); err != nil {
			return nil, err
		}
		return pinfoResp, nil
	}
}
