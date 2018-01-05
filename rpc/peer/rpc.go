package peer

import (
	"fmt"
	"github.com/Zumium/fyer/cfg"
	control_peer "github.com/Zumium/fyer/control/peer"
	pb_peer "github.com/Zumium/fyer/protos/peer"
	"github.com/xtaci/kcp-go"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
)

type rpcImpl struct {
	control_peer.DeployController
	control_peer.FetchController
}

//Start starts the RPC service
func Start() error {
	lis, err := kcp.ListenWithOptions(fmt.Sprintf("[%s]:%d", "::", cfg.Port()), nil, 10, 3)
	if err != nil {
		return err
	}
	server := grpc.NewServer(grpc.MaxSendMsgSize(cfg.MaxSendRecvMsgSize()), grpc.MaxRecvMsgSize(cfg.MaxSendRecvMsgSize()))
	pb_peer.RegisterFyerPeerServer(server, new(rpcImpl))
	return server.Serve(lis)
}
