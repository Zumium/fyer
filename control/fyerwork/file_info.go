package fyerwork

import (
	"github.com/Zumium/fyer/connectionmngr"
	pb_center "github.com/Zumium/fyer/protos/center"
	"golang.org/x/net/context"
	"time"
)

var defaultFileInfoController = new(FileInfoController)

func FileInfo(name string) (uint64, []byte, uint64, time.Time, error) {
	return defaultFileInfoController.FileInfo(name)
}

type FileInfoController struct{}

func (controller *FileInfoController) FileInfo(name string) (size uint64, hash []byte, fragCount uint64, uploadTime time.Time, err error) {
	var conn *connectionmngr.Connection
	conn, err = connectionmngr.ConnectToCenter()
	if err != nil {
		return
	}
	defer conn.Close()

	fyerCenterClient := pb_center.NewFyerCenterClient(conn.ClientConn)
	var resp *pb_center.FileInfoResponse
	resp, err = fyerCenterClient.FileInfo(context.TODO(), &pb_center.FileInfoRequest{name})
	if err != nil {
		return
	}
	return resp.GetSize(), resp.GetHash(), resp.GetFragCount(), time.Unix(resp.GetUploadTime().GetSeconds(), 0), nil
}
