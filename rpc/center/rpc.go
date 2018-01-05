package center

import (
	"fmt"
	"github.com/Zumium/fyer/cfg"
	control_center "github.com/Zumium/fyer/control/center"
	pb_center "github.com/Zumium/fyer/protos/center"
	"github.com/xtaci/kcp-go"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
)

type rpcImpl struct {
	control_center.FileRegisterController
	control_center.PeerInfoController
	control_center.FileInfoController
	control_center.FragDistributionController
	control_center.FragInfoController
}

//Start starts the center's RPC service
func Start() error {
	fmt.Printf("RPC listening on port %d\n", cfg.Port())
	lis, err := kcp.ListenWithOptions(fmt.Sprintf("[%s]:%d", "::", cfg.Port()), nil, 10, 3)
	if err != nil {
		return err
	}
	server := grpc.NewServer(grpc.MaxRecvMsgSize(cfg.MaxSendRecvMsgSize()), grpc.MaxSendMsgSize(cfg.MaxSendRecvMsgSize()))
	pb_center.RegisterFyerCenterServer(server, new(rpcImpl))
	return server.Serve(lis)
}
