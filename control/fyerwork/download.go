package fyerwork

import (
	"encoding/json"
	"fmt"
	"github.com/Zumium/fyer/alg"
	"github.com/Zumium/fyer/cfg"
	"github.com/Zumium/fyer/common"
	"github.com/Zumium/fyer/connectionmngr"
	pb_center "github.com/Zumium/fyer/protos/center"
	pb_peer "github.com/Zumium/fyer/protos/peer"
	util_peer "github.com/Zumium/fyer/util/peer"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
)

var defaultDownloadController = new(DownloadController)

//Download downloads the specified file using the default download controller
func Download(name, storePath string) error {
	return defaultDownloadController.Download(name, storePath)
}

type DownloadController struct{}

func (dc *DownloadController) getFragCount(name string) (uint64, error) {
	conn, err := connectionmngr.ConnectToCenter()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	fyerCenterClient := pb_center.NewFyerCenterClient(conn.ClientConn)
	resp, err := fyerCenterClient.FileInfo(context.TODO(), &pb_center.FileInfoRequest{Name: name})
	if err != nil {
		return 0, err
	}
	return resp.GetFragCount(), nil
}

func (dc *DownloadController) getFragInfo(name string) ([]common.Frag, error) {
	conn, err := connectionmngr.ConnectToCenter()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	fyerCenterClient := pb_center.NewFyerCenterClient(conn.ClientConn)
	resp, err := fyerCenterClient.FragInfo(context.TODO(), &pb_center.FragInfoRequest{Name: name})
	if err != nil {
		return nil, err
	}

	var frags []common.Frag
	if err := json.Unmarshal(resp.GetFrags(), &frags); err != nil {
		return nil, err
	}

	return frags, nil
}

func (dc *DownloadController) getFragDistribution(name string) ([][]string, error) {
	conn, err := connectionmngr.ConnectToCenter()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	fyerCenterClient := pb_center.NewFyerCenterClient(conn.ClientConn)
	resp, err := fyerCenterClient.FragDistribution(context.TODO(), &pb_center.FragDistributionRequest{name})
	if err != nil {
		return nil, err
	}

	result := make([][]string, 0, len(resp.Distribution))
	for i := 0; i < len(resp.Distribution); i++ {
		result = append(result, resp.Distribution[i].Peers)
	}
	return result, nil
}

func (dc *DownloadController) downloadFragAndSave(ctx context.Context, name string, frag common.Frag, peer string, file *os.File) {
	addr, err := util_peer.ResolvePeerIDByCenter(peer)
	if err != nil {
		//TODO: Deal with error
		return
	}

	conn, err := connectionmngr.ConnectTo(fmt.Sprintf("%s:%d", addr, cfg.Port()))
	if err != nil {
		//TODO: Deal with error
		return
	}
	defer conn.Close()

	fyerCenterClient := pb_peer.NewFyerPeerClient(conn.ClientConn)
	resp, err := fyerCenterClient.Fetch(context.TODO(), &pb_peer.FetchRequest{name, frag.Index})
	if err != nil {
		//TODO: Deal with error
		return
	}

	file.WriteAt(resp.Data, frag.Start)
}

func (dc *DownloadController) Download(name string, storePath string) error {
	//create target file
	fpath := filepath.Join(storePath, name)
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()
	//fetch file frag count
	//fragCount, err := dc.getFragCount(name)
	//if err != nil {
	//	return err
	//}
	//download frag info
	fragInfo, err := dc.getFragInfo(name)
	//download frag distribution table
	fragDistribution, err := dc.getFragDistribution(name)
	if err != nil {
		return err
	}
	//run download balancer algorithm
	downloadSources := alg.NewDownloadBalancer(fragDistribution).Result()
	//make fetch requests
	for i, src := range downloadSources {
		go dc.downloadFragAndSave(context.TODO(), name, fragInfo[i], src, f)
	}
	//save to file
	return nil
}
