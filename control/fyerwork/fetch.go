package fyerwork

import (
	"github.com/Zumium/fyer/filemanager"
	pb_fyerwork "github.com/Zumium/fyer/protos/fyerwork"
	"golang.org/x/net/context"
	"io"
)

type FetchController struct{}

func (f *FetchController) Fetch(ctx context.Context, in *pb_fyerwork.FetchRequest) (*pb_fyerwork.FetchResponse, error) {
	file, err := filemanager.Open(in.GetName())
	if err != nil {
		return nil, err
	}
	buf := make([]byte, in.GetRange().GetSize())

	resp := &pb_fyerwork.FetchResponse{}

	_, err = file.ReadAt(buf, in.GetRange().GetStart())
	switch err {
	case io.EOF:
		fallthrough
	case nil:
		resp.Data = buf
	default:
		return nil, err
	}
	return resp, nil
}
