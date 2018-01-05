package fyerwork

import (
	"fmt"
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

func Start(port, messageSize int) error {
	lis, err := kcp.ListenWithOptions(fmt.Sprintf("[%s]:%d", "::", port), nil, 10, 3)
	if err != nil {
		return err
	}
	server := grpc.NewServer(grpc.MaxSendMsgSize(messageSize), grpc.MaxRecvMsgSize(messageSize))
	pb_fyerwork.RegisterFyerworkServer(server, new(rpcImpl))
	err = server.Serve(lis)
	exit <- struct{}{}
	return err
}

func StartInBackground(port, messageSize int) {
	go Start(port, messageSize)
}

func Wait() {
	<-exit
}
