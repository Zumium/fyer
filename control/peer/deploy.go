package peer

import (
	"fmt"
	"github.com/Zumium/fyer/cfg"
	"github.com/Zumium/fyer/connectionmngr"
	db_peer "github.com/Zumium/fyer/db/peer"
	"github.com/Zumium/fyer/fragmngr"
	pb_center "github.com/Zumium/fyer/protos/center"
	pb_peer "github.com/Zumium/fyer/protos/peer"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"sync"
)

var (
	deploySingleton *Deploy
	deployOnce      sync.Once
)

type Deploy struct{}

//DeployInstance returns the instance of Deploy struct
//this is the only way to get the instance of Deploy
func DeployInstance() *Deploy {
	deployOnce.Do(func() {
		deploySingleton = new(Deploy)
	})
	return deploySingleton
}

func (d *Deploy) GRPCHandler() func(context.Context, *pb_peer.DeployRequest) (*google_protobuf.Empty, error) {
	return func(ctx context.Context, in *pb_peer.DeployRequest) (*google_protobuf.Empty, error) {
		//query source's IP address
		centerConn, err := connectionmngr.ConnectTo(fmt.Sprintf("%s:%d", cfg.CenterAddress(), cfg.Port()))
		if err != nil {
			return nil, err
		}
		defer centerConn.Close()
		centerClient := pb_center.NewFyerCenterClient(centerConn.ClientConn)
		peerInfoResp, err := centerClient.PeerInfo(context.TODO(), &pb_center.PeerInfoRequest{PeerId: in.GetSource()})
		if err != nil {
			return nil, err
		}

		conn, err := connectionmngr.ConnectTo(fmt.Sprintf("%s:%d", peerInfoResp.GetAddress(), cfg.Port()))
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		peerClient := pb_peer.NewFyerPeerClient(conn.ClientConn)
		resp, err := peerClient.Fetch(context.TODO(), &pb_peer.FetchRequest{Name: in.GetName(), FragIndex: in.GetFragIndex()})
		if err != nil {
			return nil, err
		}

		fileAdapter, err := fragmngr.FMInstance().Open(in.GetName())
		if err != nil {
			return nil, err
		}
		defer fileAdapter.Close()
		if err := fileAdapter.Write(in.GetFragIndex(), resp.GetData()); err != nil {
			return nil, err
		}

		fileDB := db_peer.ToFile(in.GetName())
		exist := fileDB.Has()
		if fileDB.Err() != nil {
			return nil, fileDB.Err()
		}
		if exist {
			storedFrags := fileDB.StoredFrags()
			if fileDB.Err() != nil {
				return nil, fileDB.Err()
			}
			storedFrags.Add(in.GetFragIndex())
			if err := fileDB.Edit().SetStoredFrags(storedFrags).Done(); err != nil {
				return nil, err
			}
		} else {
			if err := fileDB.Edit().SetStoredFrags(&db_peer.StoredFrags{Numbers: []uint64{in.GetFragIndex()}}).Done(); err != nil {
				return nil, err
			}
		}

		return new(google_protobuf.Empty), nil
	}
}
