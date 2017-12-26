package center

import (
	"github.com/Zumium/fyer/alg"
	"github.com/Zumium/fyer/cfg"
	common_peer "github.com/Zumium/fyer/common/peer"
	"github.com/Zumium/fyer/connectionmngr"
	db_center "github.com/Zumium/fyer/db/center"
	"github.com/Zumium/fyer/merkle"
	pb_center "github.com/Zumium/fyer/protos/center"
	pb_peer "github.com/Zumium/fyer/protos/peer"
	util_center "github.com/Zumium/fyer/util/center"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"time"
)

//FileRegisterStoreFileInfo is used to store file info into database
//it represents an abstract db operating process
// type FileRegisterStoreFileInfo interface {
// 	Store(name string, size uint64, hash []byte, fragCount uint64, mtree *merkle.MTree)
// }

//FileRegisterController -- file registering process controller
type FileRegisterController struct{}

//storeToDB stores file meta data to db
func (fr *FileRegisterController) storeToDB(in *pb_center.RegisterRequest) error {
	dbHandler, err := db_center.ToFileMeta(in.Name)
	if err != nil {
		return err
	}
	mtree, err := merkle.Unmarshal(in.MerkleTree)
	if err != nil {
		return err
	}
	editor := dbHandler.Edit()
	editor.SetSize(in.Size).SetHash(in.Hash).SetFragCount(in.FragCount).SetMerkleTree(mtree).SetUploadTime(time.Now())
	if err := editor.Done(); err != nil {
		return err
	}
	return nil
}

//dispatchFrags makes deploying requests to peers
func (fr *FileRegisterController) dispatchFrags(in *pb_center.RegisterRequest) ([]common_peer.Frag, map[uint64][]string, error) {
	peers, err := db_center.AllPeers()
	if err != nil {
		return nil, nil, err
	}
	peerIDs := make([]string, 0, len(peers))
	for _, p := range peers {
		peerIDs = append(peerIDs, p.PeerID())
	}

	fragCutter := &alg.FragCutter{Size: in.GetSize(), FragSize: cfg.FragSize()}
	frags := fragCutter.Cut()
	dispatches := alg.NewDispatchAlg(peerIDs, uint64(len(frags)), cfg.Replica()).Dispatch()

	return frags, dispatches, nil
}

//makeDeploys make RPC deploying calls to peers
func (fr *FileRegisterController) makeDeploys(frags []common_peer.Frag, dispatches map[uint64][]string, in *pb_center.RegisterRequest) error {
	for i, dispatch := range dispatches { // i -- frag index, dispatch -- array of peers' id on whom the frag will be deployed
		for j, peerID := range dispatch {
			pbFrag := common_peer.FragCommonToPb(frags[i])

			address, err := util_center.ResolvePeerID(peerID)
			if err != nil {
				return err
			}
			conn, err := connectionmngr.ConnectTo(address)
			if err != nil {
				return err
			}
			peerClient := pb_peer.NewFyerPeerClient(conn.ClientConn)
			if j == 0 {
				//deploy frag from client to the first peer
				_, err = peerClient.Deploy(context.TODO(), &pb_peer.DeployRequest{in.Name, pbFrag, in.GetSource(), pb_peer.DeployRequest_CLIENT})
				if err != nil {
					conn.Close()
					return err
				}
			} else {
				//deploy frag from last peer to next one
				_, err = peerClient.Deploy(context.TODO(), &pb_peer.DeployRequest{in.Name, pbFrag, dispatch[j-1], pb_peer.DeployRequest_PEER})
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
	if err := fr.storeToDB(in); err != nil {
		return nil, err
	}
	frags, dispatches, err := fr.dispatchFrags(in)
	if err != nil {
		return nil, err
	}
	if err = fr.makeDeploys(frags, dispatches, in); err != nil {
		return nil, err
	}

	return new(google_protobuf.Empty), nil
}
