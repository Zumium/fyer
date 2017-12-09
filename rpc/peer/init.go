package peer

import (
	"fmt"
	"github.com/Zumium/fyer/cfg"
	control_peer "github.com/Zumium/fyer/control/peer"
	pb_peer "github.com/Zumium/fyer/protos/peer"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/xtaci/kcp-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type rpcImpl struct{}

func (i *rpcImpl) Deploy(ctx context.Context, in *pb_peer.DeployRequest) (*google_protobuf.Empty, error) {
	return control_peer.DeployInstance().GRPCHandler()(ctx, in)
}

func (i *rpcImpl) Fetch(ctx context.Context, in *pb_peer.FetchRequest) (*pb_peer.FetchResponse, error) {
	return control_peer.FetchInstance().GRPCHandler()(ctx, in)
}

//Start starts the RPC service
func Start() error {
	lis, err := kcp.ListenWithOptions(fmt.Sprintf("[%s]:%d", "::", cfg.Port()), nil, 10, 3)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb_peer.RegisterFyerPeerServer(server, new(rpcImpl))
	return server.Serve(lis)
}
