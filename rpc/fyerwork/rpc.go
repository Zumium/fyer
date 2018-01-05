package fyerwork

import (
	"fmt"
	"github.com/Zumium/fyer/cfg"
	"github.com/Zumium/fyer/control/fyerwork"
	pb_fyerwork "github.com/Zumium/fyer/protos/fyerwork"
	"github.com/xtaci/kcp-go"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
)

var exit = make(chan struct{}, 1)

type rpcImpl struct {
	fyerwork.FetchController
}

func Start() error {
	fmt.Printf("RPC listening on port %d\n", cfg.Port())
	lis, err := kcp.ListenWithOptions(fmt.Sprintf("[%s]:%d", "::", cfg.Port()), nil, 10, 3)
	if err != nil {
		return err
	}
	server := grpc.NewServer(grpc.MaxSendMsgSize(cfg.MaxSendRecvMsgSize()), grpc.MaxRecvMsgSize(cfg.MaxSendRecvMsgSize()))
	pb_fyerwork.RegisterFyerworkServer(server, new(rpcImpl))
	err = server.Serve(lis)
	exit <- struct{}{}
	return err
}

func StartInBackground() {
	go Start()
}

func Wait() {
	<-exit
}
