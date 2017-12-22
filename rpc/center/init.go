package center

import (
	"fmt"
	"github.com/Zumium/fyer/cfg"
	control_center "github.com/Zumium/fyer/control/center"
	pb_center "github.com/Zumium/fyer/protos/center"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"net"
)

type rpcImpl struct{}

func (i *rpcImpl) Register(ctx context.Context, in *pb_center.RegisterRequest) (*google_protobuf.Empty, error) {
	return control_center.FileRegisterInstance().GRPCHandler()(ctx, in)
}

func (i *rpcImpl) FileInfo(ctx context.Context, in *pb_center.FileInfoRequest) (*pb_center.FileInfoResponse, error) {
	return control_center.FileInfoInstance().GRPCHandler()(ctx, in)
}

func (i *rpcImpl) PeerInfo(ctx context.Context, in *pb_center.PeerInfoRequest) (*pb_center.PeerInfoResponse, error) {
	return control_center.PeerInfoInstance().GRPCHandler()(ctx, in)
}

//Start starts the center's RPC service
func Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("[%s]:%d", "::", cfg.Port()))
	if err != nil {
		return err
	}
	server := grpc.NewServer(grpc.MaxRecvMsgSize(cfg.MaxSendRecvMsgSize()), grpc.MaxSendMsgSize(cfg.MaxSendRecvMsgSize()))
	pb_center.RegisterFyerCenterServer(server, new(rpcImpl))
	return server.Serve(lis)
}
