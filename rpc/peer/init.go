package peer

import (
	"fmt"
	"github.com/Zumium/fyer/cfg"
	control_peer "github.com/Zumium/fyer/control/peer"
	pb_peer "github.com/Zumium/fyer/protos/peer"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"net"
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
	lis, err := net.Listen("tcp", fmt.Sprintf("[%s]:%d", "::", cfg.Port()))
	if err != nil {
		return err
	}
	server := grpc.NewServer(grpc.MaxSendMsgSize(cfg.MaxSendRecvMsgSize()), grpc.MaxRecvMsgSize(cfg.MaxSendRecvMsgSize()))
	pb_peer.RegisterFyerPeerServer(server, new(rpcImpl))
	return server.Serve(lis)
}
