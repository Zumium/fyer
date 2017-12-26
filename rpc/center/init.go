package center

import (
	"fmt"
	"github.com/Zumium/fyer/cfg"
	control_center "github.com/Zumium/fyer/control/center"
	pb_center "github.com/Zumium/fyer/protos/center"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"net"
)

type rpcImpl struct {
	control_center.FileRegisterController
	control_center.PeerInfoController
	control_center.FileInfoController
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
