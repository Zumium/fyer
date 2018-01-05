package peer

import (
	"fmt"
	"github.com/Zumium/fyer/cfg"
	"github.com/Zumium/fyer/common"
	"github.com/Zumium/fyer/connectionmngr"
	db_peer "github.com/Zumium/fyer/db/peer"
	"github.com/Zumium/fyer/fragmngr"
	pb_fyerwork "github.com/Zumium/fyer/protos/fyerwork"
	pb_peer "github.com/Zumium/fyer/protos/peer"
	util_peer "github.com/Zumium/fyer/util/peer"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

type DeployController struct{}

func (d *DeployController) fetchFragDataFromFyerwork(in *pb_peer.DeployRequest) ([]byte, error) {
	fmt.Printf("Establishing connection to client %s\n", in.GetSrc())
	conn, err := connectionmngr.ConnectTo(in.GetSrc())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	frag := common.MustUnmarshalJsonToFrag(in.GetFrag())
	fyerworkClient := pb_fyerwork.NewFyerworkClient(conn.ClientConn)
	resp, err := fyerworkClient.Fetch(context.TODO(), &pb_fyerwork.FetchRequest{
		Name: in.GetName(),
		Range: &pb_fyerwork.FetchRequest_Range{
			Start: frag.Start,
			Size:  frag.Size,
		},
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("Retrieved %d bytes of data\n", len(resp.GetData()))
	return resp.GetData(), nil
}

func (d *DeployController) fetchFragDataFromPeer(address string, in *pb_peer.DeployRequest) ([]byte, error) {
	fmt.Printf("Establishing connection to peer %s\n", address)
	//get the connection to the source
	conn, err := connectionmngr.ConnectTo(fmt.Sprintf("%s:%d", address, cfg.Port()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	//retreive the frag data
	peerClient := pb_peer.NewFyerPeerClient(conn.ClientConn)
	frag := common.MustUnmarshalJsonToFrag(in.GetFrag())
	resp, err := peerClient.Fetch(context.TODO(), &pb_peer.FetchRequest{Name: in.GetName(), FragIndex: frag.Index})
	if err != nil {
		return nil, err
	}

	fmt.Printf("Retrieved %d bytes of data\n", len(resp.GetData()))
	return resp.GetData(), nil
}

func (d *DeployController) writeFragData(data []byte, in *pb_peer.DeployRequest) error {
	fileAdapter, err := fragmngr.FMInstance().Open(in.GetName())
	if err != nil {
		return err
	}
	defer fileAdapter.Close()
	frag := common.MustUnmarshalJsonToFrag(in.GetFrag())
	fmt.Printf("Wring %d bytes of data to storage: %s\n", len(data), frag.String())
	if err := fileAdapter.Write(frag, data); err != nil {
		return err
	}
	return nil
}

func (d *DeployController) addDBRecord(in *pb_peer.DeployRequest) error {
	fileDB := db_peer.ToFile(in.GetName())
	exist := fileDB.Has()
	if fileDB.Err() != nil {
		return fileDB.Err()
	}
	if exist {
		fmt.Printf("file record %s already exists\n", in.GetName())
		storedFrags := fileDB.StoredFrags()
		if fileDB.Err() != nil {
			return fileDB.Err()
		}
		storedFrags.AddFrag(common.MustUnmarshalJsonToFrag(in.GetFrag()))
		if err := fileDB.Edit().SetStoredFrags(storedFrags).Done(); err != nil {
			return err
		}
	} else {
		fmt.Printf("file record %s doesnt exist, create a new record\n", in.GetName())
		if err := fileDB.Edit().SetStoredFrags(&db_peer.StoredFrags{Frags: []common.Frag{common.MustUnmarshalJsonToFrag(in.GetFrag())}}).Done(); err != nil {
			return err
		}
	}
	return nil
}

func (d *DeployController) Deploy(ctx context.Context, in *pb_peer.DeployRequest) (*google_protobuf.Empty, error) {
	var err error
	var data []byte

	fmt.Println("Received a new deploy request")
	fmt.Printf("New deploy request: %s\n", in.String())

	switch in.GetSrcType() {
	case pb_peer.DeployRequest_CLIENT:
		//download the specified frag data from client, then put it in "data"
		fmt.Println("Fetch frag from client")
		data, err = d.fetchFragDataFromFyerwork(in)
	case pb_peer.DeployRequest_PEER:
		fmt.Println("Fetch frag from peer")
		address, err := util_peer.ResolvePeerIDByCenter(in.GetSrc())
		if err != nil {
			return nil, err
		}
		fmt.Printf("source: %s - address: %s\n", in.GetSrc(), address)
		data, err = d.fetchFragDataFromPeer(address, in)
	}
	if err != nil {
		return nil, err
	}

	//store frag data
	fmt.Println("Storing frag data")
	if err := d.writeFragData(data, in); err != nil {
		return nil, err
	}
	//add db record
	fmt.Println("Saving frag record to database")
	if err := d.addDBRecord(in); err != nil {
		return nil, err
	}

	return new(google_protobuf.Empty), nil
}
