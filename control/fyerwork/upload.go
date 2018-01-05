package fyerwork

import (
	"fmt"
	"github.com/Zumium/fyer/cfg"
	"github.com/Zumium/fyer/connectionmngr"
	pb_center "github.com/Zumium/fyer/protos/center"
	"golang.org/x/net/context"
)

var defaultUploadController = new(UploadController)

func UploadFile(name string, size uint64, hash []byte, laddr string) error {
	return defaultUploadController.Upload(name, size, hash, fmt.Sprintf("%s:%d", laddr, cfg.Port()))
}

type UploadController struct{}

func (controller *UploadController) Upload(name string, size uint64, hash []byte, laddr string) error {
	conn, err := connectionmngr.ConnectToCenter()
	if err != nil {
		return err
	}
	defer conn.Close()

	fyerCenterClient := pb_center.NewFyerCenterClient(conn.ClientConn)
	_, err = fyerCenterClient.Register(context.TODO(), &pb_center.RegisterRequest{name, size, hash, laddr})
	return err
}
