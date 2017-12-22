package peer

import (
	"fmt"
	"github.com/Zumium/fyer/cfg"
	common_peer "github.com/Zumium/fyer/common/peer"
	"github.com/Zumium/fyer/connectionmngr"
	db_peer "github.com/Zumium/fyer/db/peer"
	"github.com/Zumium/fyer/fragmngr"
	pb_fyerwork "github.com/Zumium/fyer/protos/fyerwork"
	pb_peer "github.com/Zumium/fyer/protos/peer"
	util_peer "github.com/Zumium/fyer/util/peer"
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

func (d *Deploy) fetchFragDataFromFyerwork(in *pb_peer.DeployRequest) ([]byte, error) {
	conn, err := connectionmngr.ConnectTo(in.GetSrc())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	fyerworkClient := pb_fyerwork.NewFyerworkClient(conn.ClientConn)
	resp, err := fyerworkClient.Fetch(context.TODO(), &pb_fyerwork.FetchRequest{
		Name: in.GetName(),
		Range: &pb_fyerwork.FetchRequest_Range{
			Start: in.GetFrag().Start,
			Size:  in.GetFrag().Size,
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.GetData(), nil
}

func (d *Deploy) fetchFragDataFromPeer(address string, in *pb_peer.DeployRequest) ([]byte, error) {
	//get the connection to the source
	conn, err := connectionmngr.ConnectTo(fmt.Sprintf("%s:%d", address, cfg.Port()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	//retreive the frag data
	peerClient := pb_peer.NewFyerPeerClient(conn.ClientConn)
	resp, err := peerClient.Fetch(context.TODO(), &pb_peer.FetchRequest{Name: in.GetName(), FragIndex: in.GetFrag().Index})
	if err != nil {
		return nil, err
	}
	return resp.GetData(), nil
}

func (d *Deploy) writeFragData(data []byte, in *pb_peer.DeployRequest) error {
	fileAdapter, err := fragmngr.FMInstance().Open(in.GetName())
	if err != nil {
		return err
	}
	defer fileAdapter.Close()
	if err := fileAdapter.Write(common_peer.FragPbToCommon(in.GetFrag()), data); err != nil {
		return err
	}
	return nil
}

func (d *Deploy) addDBRecord(in *pb_peer.DeployRequest) error {
	fileDB := db_peer.ToFile(in.GetName())
	exist := fileDB.Has()
	if fileDB.Err() != nil {
		return fileDB.Err()
	}
	if exist {
		storedFrags := fileDB.StoredFrags()
		if fileDB.Err() != nil {
			return fileDB.Err()
		}
		storedFrags.AddFrag(common_peer.FragPbToCommon(in.GetFrag()))
		if err := fileDB.Edit().SetStoredFrags(storedFrags).Done(); err != nil {
			return err
		}
	} else {
		if err := fileDB.Edit().SetStoredFrags(&db_peer.StoredFrags{Frags: []common_peer.Frag{common_peer.FragPbToCommon(in.GetFrag())}}).Done(); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deploy) GRPCHandler() func(context.Context, *pb_peer.DeployRequest) (*google_protobuf.Empty, error) {
	return func(ctx context.Context, in *pb_peer.DeployRequest) (*google_protobuf.Empty, error) {
		var err error
		var data []byte

		switch in.GetSrcType() {
		case pb_peer.DeployRequest_CLIENT:
			//download the specified frag data from client, then put it in "data"
			data, err = d.fetchFragDataFromFyerwork(in)
		case pb_peer.DeployRequest_PEER:
			address, err := util_peer.ResolvePeerIDByCenter(in.GetSrc())
			if err != nil {
				return nil, err
			}
			data, err = d.fetchFragDataFromPeer(address, in)
		}
		if err != nil {
			return nil, err
		}

		//store frag data
		if err := d.writeFragData(data, in); err != nil {
			return nil, err
		}
		//add db record
		if err := d.addDBRecord(in); err != nil {
			return nil, err
		}

		return new(google_protobuf.Empty), nil
	}
}
