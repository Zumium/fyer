package peer

import (
	"github.com/Zumium/fyer/connectionmngr"
	pb_center "github.com/Zumium/fyer/protos/center"
	"golang.org/x/net/context"
)

func ResolvePeerIDByCenter(peerID string) (string, error) {
	centerConn, err := connectionmngr.ConnectToCenter()
	if err != nil {
		return peerID, err
	}
	defer centerConn.Close()
	centerClient := pb_center.NewFyerCenterClient(centerConn.ClientConn)
	peerInfoResp, err := centerClient.PeerInfo(context.TODO(), &pb_center.PeerInfoRequest{PeerId: peerID})
	if err != nil {
		return peerID, err
	}
	return peerInfoResp.GetAddress(), nil
}
