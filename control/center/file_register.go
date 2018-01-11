package center

import (
	"github.com/Zumium/fyer/alg"
	"github.com/Zumium/fyer/cfg"
	"github.com/Zumium/fyer/common"
	"github.com/Zumium/fyer/connectionmngr"
	db_center "github.com/Zumium/fyer/db/center"
	"fmt"
	pb_center "github.com/Zumium/fyer/protos/center"
	pb_peer "github.com/Zumium/fyer/protos/peer"
	util_center "github.com/Zumium/fyer/util/center"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"time"
)

//FileRegisterController -- file registering process controller
type FileRegisterController struct{}

//storeToDB stores file meta data to db
func (fr *FileRegisterController) storeToDB(in *pb_center.RegisterRequest) error {
	fmt.Printf("Storing file meta data: %s\n", in.String())
	dbHandler, err := db_center.ToFileMeta(in.GetName())
	if err != nil {
		return err
	}
	return dbHandler.Edit().SetSize(in.GetSize()).SetHash(in.GetHash()).SetUploadTime(time.Now()).Done()
}

//dispatchFrags makes deploying requests to peers
func (fr *FileRegisterController) dispatchFrags(in *pb_center.RegisterRequest) ([]common.Frag, [][]string, error) {
	peers, err := db_center.AllPeers()
	if err != nil {
		return nil, nil, err
	}
	peerIDs := make([]string, 0, len(peers))
	for _, p := range peers {
		peerIDs = append(peerIDs, p.PeerID())
	}
	fmt.Printf("current peer ids: %s\n", peerIDs)

	//fragCutter := &alg.FragCutter{Size: in.GetSize(), FragSize: cfg.FragSize()}
	fragCutter := alg.NewFragCutter(in.GetSize(), cfg.FragSize())
	frags := fragCutter.Cut()
	for _, f := range frags {
		fmt.Println(f.String())
	}
	dispatches := alg.NewDispatchAlg(peerIDs, uint64(len(frags)), cfg.Replica()).Dispatch()
	fmt.Printf("dispatch result: %v\n", dispatches)

	return frags, dispatches, nil
}

//saveFragCount saves numbers of frag to db
func (fr *FileRegisterController) saveFragCount(name string, fragCount uint64) error {
	dbHandler, err := db_center.ToFileMeta(name)
	if err != nil {
		return err
	}
	return dbHandler.Edit().SetFragCount(fragCount).Done()
}

//saveFragInfo saves frags' details and peer lists
func (fr *FileRegisterController) saveFragInfo(name string, frags []common.Frag, peerLists [][]string) error {
	handler, err := db_center.ToFragFile(name)
	if err != nil {
		return err
	}
	return handler.Edit().SetFrags(frags).SetPeerList(peerLists).Done()
}

//makeDeploys make RPC deploying calls to peers
func (fr *FileRegisterController) makeDeploys(frags []common.Frag, dispatches [][]string, in *pb_center.RegisterRequest) error {
	for i, dispatch := range dispatches { // i -- frag index, dispatch -- array of peers' id on whom the frag will be deployed
		for j, peerID := range dispatch {
			fmt.Printf("Deploying frag %d to peer %s\n", i, peerID)
			bFrag := common.MustMarshalFragAsJson(frags[i])

			address, err := util_center.ResolvePeerID(peerID)
			if err != nil {
				return err
			}
			fmt.Printf("peer %s address: %s\n", peerID, address)

			conn, err := connectionmngr.ConnectToWithDefaultPort(address)
			if err != nil {
				return err
			}
			peerClient := pb_peer.NewFyerPeerClient(conn.ClientConn)
			if j == 0 {
				fmt.Println("Fetch frag from client")
				//deploy frag from client to the first peer
				_, err = peerClient.Deploy(context.TODO(), &pb_peer.DeployRequest{in.Name, bFrag, in.GetSource(), pb_peer.DeployRequest_CLIENT})
				if err != nil {
					conn.Close()
					return err
				}
			} else {
				fmt.Println("Fetch frag from previous peer")
				//deploy frag from last peer to next one
				_, err = peerClient.Deploy(context.TODO(), &pb_peer.DeployRequest{in.Name, bFrag, dispatch[j-1], pb_peer.DeployRequest_PEER})
				if err != nil {
					conn.Close()
					return err
				}
			}
			conn.Close()
		}
	}
	return nil
}

//handles GRPC request
func (fr *FileRegisterController) Register(ctx context.Context, in *pb_center.RegisterRequest) (*google_protobuf.Empty, error) {
	fmt.Printf("Received a register request: %s\n", in.String())
	fmt.Println("Storing file meta data in database")
	if err := fr.storeToDB(in); err != nil {
		return nil, err
	}

	fmt.Println("Cutting fragments")
	frags, dispatches, err := fr.dispatchFrags(in)
	if err != nil {
		return nil, err
	}

	fmt.Println("Saving frag count")
	if err := fr.saveFragCount(in.GetName(), uint64(len(frags))); err != nil {
		return nil, err
	}

	fmt.Println("Saving frags' info")
	if err := fr.saveFragInfo(in.GetName(), frags, dispatches); err != nil {
		return nil, err
	}

	fmt.Println("Making deploying requests")
	if err = fr.makeDeploys(frags, dispatches, in); err != nil {
		return nil, err
	}

	return new(google_protobuf.Empty), nil
}
