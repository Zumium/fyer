package peer

import (
	"errors"
	"github.com/Zumium/fyer/common"
	db_peer "github.com/Zumium/fyer/db/peer"
	"github.com/Zumium/fyer/fragmngr"
	pb_peer "github.com/Zumium/fyer/protos/peer"
	"golang.org/x/net/context"
	"fmt"
)

var (
	ErrFileNotFound = errors.New("File not found")
	ErrFragNotFound = errors.New("Frag not found")
)

type FetchController struct{}

func (ftch *FetchController) checkFileAndFragExist(in *pb_peer.FetchRequest) error {
	//check whether this peer actually is holding the specified fragment
	f := db_peer.ToFile(in.GetName())
	exist := f.Has()
	if f.Err() != nil {
		return f.Err()
	}
	if !exist {
		return ErrFileNotFound
	}
	//check whether the specified frag exists locally
	storedFrags := f.StoredFrags()
	if f.Err() != nil {
		return f.Err()
	}
	if !storedFrags.Has(in.GetFragIndex()) {
		return ErrFragNotFound
	}
	return nil
}

func (ftch *FetchController) readLocalFrag(in *pb_peer.FetchRequest) common.Frag {
	sf := db_peer.ToFile(in.GetName()).StoredFrags()
	frag, _ := sf.Find(in.GetFragIndex())
	return frag
}

func (ftch *FetchController) Fetch(ctx context.Context, in *pb_peer.FetchRequest) (*pb_peer.FetchResponse, error) {
	//check whether file and frag exist
	fmt.Printf("New fetch request: %s\n",in.String())
	fmt.Println(in.String())
	if err := ftch.checkFileAndFragExist(in); err != nil {
		return nil, err
	}

	//read the specified data out and return the response
	fmt.Println("Reading file frag")
	fileApater, err := fragmngr.FMInstance().Open(in.GetName())
	if err != nil {
		return nil, err
	}
	defer fileApater.Close()

	//read frag data
	d, err := fileApater.Read(ftch.readLocalFrag(in))
	if err != nil {
		return nil, err
	}
	fmt.Printf("Read %d bytes out\n", len(d))

	return &pb_peer.FetchResponse{Data: d}, nil
}
