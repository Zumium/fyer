package merkle

import (
	"bytes"
	"hash"
	"math"

	pb_merkle "github.com/Zumium/fyer/protos/util"
	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/sha3"
)

//MTree - Merkle Tree
type MTree struct {
	nodes  [][]byte // COMPLETE BINARY TREE
	hasher hash.Hash
}

//NewTree creates an empty merkle tree with given capability
// func NewTree(capa uint64) *MTree {
// 	if capa < 2 {
// 		capa = 2
// 	}
// 	return &MTree{nodes: make([][]byte, 0, capa), hasher:sha3.New256()}
// }

//NewTreeFrom builds a merkle tree from given bytes fragment and returns it
func NewTreeFrom(d [][]byte) *MTree {
	l := uint64(len(d))
	if l < 1 {
		return nil
	}

	level := int(math.Ceil(math.Log2(float64(l)) + 1)) //CALCULATE THE DEPTH
	size := uint64(math.Pow(2, float64(level-1))-1) + l
	mtree := &MTree{nodes: make([][]byte, size), hasher: sha3.New256()}

	lo, hi := mtree.leavesIndex()
	for i, j := lo, uint64(0); i < hi && j < l; i, j = i+1, j+1 {
		mtree.hasher.Reset()
		mtree.hasher.Write(d[j])
		mtree.nodes[i] = mtree.hasher.Sum(nil)
	}

	for k := level - 1; k > 0; k-- {
		begin, end := uint64(math.Pow(2, float64(k-1))-1), uint64(math.Pow(2, float64(k))-1)
		for i := begin; i < end; i++ {
			var left, right []byte = nil, nil
			if leftIdx := 2*i + 1; leftIdx < size {
				left = mtree.nodes[leftIdx]
			}
			if rightIdx := 2*i + 2; rightIdx < size {
				right = mtree.nodes[2*i+2]
			}
			if left == nil && right == nil {
				mtree.nodes[i] = nil
			} else if left != nil && right != nil {
				mtree.hasher.Reset()
				mtree.hasher.Write(left)
				mtree.hasher.Write(right)
				mtree.nodes[i] = mtree.hasher.Sum(nil)
			} else if left != nil {
				mtree.hasher.Reset()
				mtree.hasher.Write(left)
				mtree.nodes[i] = mtree.hasher.Sum(nil)
			} else if right != nil {
				mtree.hasher.Reset()
				mtree.hasher.Write(right)
				mtree.nodes[i] = mtree.hasher.Sum(nil)
			}
		}
	}

	return mtree //TO BE REMOVED, REMEMBER!
}

//Marshal serialize merkle tree to bytes
func Marshal(tree *MTree) ([]byte, error) {
	merklePack := &pb_merkle.MerklePack{
		Size:  tree.Size(),
		Items: make([]*pb_merkle.MerkleItem, 0, tree.Size()),
	}
	for _, item := range tree.nodes {
		if item == nil {
			merklePack.Items = append(merklePack.Items, &pb_merkle.MerkleItem{
				Type: pb_merkle.MerkleItem_EMPTY,
			})
		} else {
			merklePack.Items = append(merklePack.Items, &pb_merkle.MerkleItem{
				Type: pb_merkle.MerkleItem_DATA,
				Data: item,
			})
		}
	}
	b, err := proto.Marshal(merklePack)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//Unmarshal deserialize bytes to merkle tree
func Unmarshal(b []byte) (*MTree, error) {
	merklePack := &pb_merkle.MerklePack{}
	if err := proto.Unmarshal(b, merklePack); err != nil {
		return nil, err
	}
	mtree := &MTree{
		nodes:  make([][]byte, merklePack.GetSize()),
		hasher: sha3.New256(),
	}

	items := merklePack.GetItems()
	for i := uint64(0); i < merklePack.GetSize(); i++ {
		switch items[i].GetType() {
		case pb_merkle.MerkleItem_DATA:
			mtree.nodes[i] = items[i].GetData()
		case pb_merkle.MerkleItem_EMPTY:
			mtree.nodes[i] = nil
		}
	}

	return mtree, nil
}

//Size returns the size of the merkle tree
func (t *MTree) Size() uint64 {
	return uint64(len(t.nodes))
}

//Root returns the merkle tree root
func (t *MTree) Root() []byte {
	/* SHOULD NEVER HAPPEN */
	if t.Size() < 1 {
		return nil
	}
	/* SHOULD NEVER HAPPEN */
	return t.nodes[0]
}

//CopyGet returns an copy data at the given position
func (t *MTree) CopyGet(idx uint64) []byte {
	if idx >= t.Size() {
		return nil
	}
	src := t.nodes[idx]
	rtn := make([]byte, len(src))
	copy(rtn, src)
	return rtn
}

//Leaves returns the leave nodes as a slice
func (t *MTree) Leaves() [][]byte {
	if t.Size() < 1 {
		return nil
	}
	lo, _ := t.leavesIndex()
	return t.nodes[lo:]
}

//SelfValid checks if the merkle tree itself is legal
func (t *MTree) SelfValid() bool {
	size := uint64(t.Size())
	for i := uint64(0); i < uint64(math.Pow(2, float64(t.depth()-1))-1); i++ {
		if leftIdx, rightIdx := 2*i+1, 2*i+2; t.nodes[i] == nil && !(leftIdx >= size && rightIdx >= size) {
			return false
		} else if l, r := leftIdx < size, rightIdx < size; l && r {
			left, right := t.nodes[2*i+1], t.nodes[2*i+2]
			t.hasher.Reset()
			t.hasher.Write(left)
			t.hasher.Write(right)
			if !bytes.Equal(t.hasher.Sum(nil), t.nodes[i]) {
				return false
			}
		} else if l {
			t.hasher.Reset()
			t.hasher.Write(t.nodes[leftIdx])
			if !bytes.Equal(t.hasher.Sum(nil), t.nodes[i]) {
				return false
			}
		} else if r {
			t.hasher.Reset()
			t.hasher.Write(t.nodes[rightIdx])
			if !bytes.Equal(t.hasher.Sum(nil), t.nodes[i]) {
				return false
			}
		}
	}
	return true
}

//FullValid contains self valid and origin data valid
//TO BE IMPLEMENTED

func (t *MTree) leavesIndex() (uint64, uint64) {
	size := t.Size()
	level := t.depth()
	lo := uint64(math.Pow(2, float64(level-1)) - 1)
	return lo, size
}

func (t *MTree) depth() int {
	return int(math.Ceil(math.Log2(float64(t.Size() + 1))))
}
