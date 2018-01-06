package fyerwork

import (
	"github.com/Zumium/fyer/connectionmngr"
	pb_center "github.com/Zumium/fyer/protos/center"
	google_protobuf_empty "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

func Files() ([]string, error) {
	return defaultFilesController.Files()
}

var defaultFilesController = new(FilesController)

type FilesController struct{}

func (controller *FilesController) Files() ([]string, error) {
	conn, err := connectionmngr.ConnectToCenter()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb_center.NewFyerCenterClient(conn.ClientConn)
	resp, err := client.Files(context.TODO(), new(google_protobuf_empty.Empty))
	if err != nil {
		return nil, err
	}
	return resp.GetFiles(), nil
}
