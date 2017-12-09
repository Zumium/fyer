package peer

import (
	"errors"
	db_peer "github.com/Zumium/fyer/db/peer"
	"github.com/Zumium/fyer/fragmngr"
	pb_peer "github.com/Zumium/fyer/protos/peer"
	"golang.org/x/net/context"
	"sync"
)

var (
	ErrFileNotFound = errors.New("File not found")
	ErrFragNotFound = errors.New("Frag not found")
)

var (
	fetchSingleton *Fetch
	fetchOnce      sync.Once
)

type Fetch struct{}

func FetchInstance() *Fetch {
	fetchOnce.Do(func() {
		fetchSingleton = new(Fetch)
	})
	return fetchSingleton
}

func (f *Fetch) GRPCHandler() func(context.Context, *pb_peer.FetchRequest) (*pb_peer.FetchResponse, error) {
	return func(ctx context.Context, in *pb_peer.FetchRequest) (*pb_peer.FetchResponse, error) {
		//check whether this peer actually is holding the specified fragment
		f := db_peer.ToFile(in.GetName())
		exist := f.Has()
		if f.Err() != nil {
			return nil, f.Err()
		}
		if !exist {
			return nil, ErrFileNotFound
		}
		//check whether the specified frag exists locally
		storedFrags := f.StoredFrags()
		if f.Err() != nil {
			return nil, f.Err()
		}
		if !storedFrags.Has(in.GetFragIndex()) {
			return nil, ErrFragNotFound
		}
		//read the specified data out and return the response
		fileApater, err := fragmngr.FMInstance().Open(in.GetName())
		if err != nil {
			return nil, err
		}
		defer fileApater.Close()

		d, err := fileApater.Read(in.GetFragIndex())
		if err != nil {
			return nil, err
		}

		return &pb_peer.FetchResponse{Data: d}, nil
	}
}
