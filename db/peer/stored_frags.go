package peer

import (
	common_peer "github.com/Zumium/fyer/common/peer"
	"github.com/coreos/etcd/store"
)

//StoredFrags represents fragment index numbers which has been stored locally
// type StoredFrags []uint64
type StoredFrags struct {
	Frags []common_peer.Frag `json:"frags"`
}

//NewEmptyStoredFrags creates an new empty StoredFrags struct
func NewEmptyStoredFrags() *StoredFrags {
	return new(StoredFrags)
}

//Has returns true if the n-th fragment has been stored locally
func (sf *StoredFrags) Has(n uint64) bool {
	l := len(sf.Frags)
	switch l {
	case 0:
		return false
	case 1:
		return sf.Frags[0].Index == n
	}
	if n < sf.Frags[0].Index || n > sf.Frags[l-1].Index {
		return false
	}
	//二分查找
	var (
		begin uint64
		end   = uint64(l - 1)
		mid   uint64
	)
	for begin < end {
		if sf.Frags[begin].Index == n || sf.Frags[end].Index == n {
			return true
		}
		mid = (begin + end) / 2
		if n < sf.Frags[mid].Index {
			end = mid
		} else {
			begin = mid + 1
		}
	}
	return false
}

//AddFrag inserts a new Frag while keep it sorted, returns the result whether
//the frag has been successfully inserted
func (sf *StoredFrags) AddFrag(frag common_peer.Frag) bool {
	for i := 0; i < len(sf.Frags); i++ {
		if sf.Frags[i].Index < frag.Index {
			continue
		} else if sf.Frags[i].Index == frag.Index {
			return false
		} else {
			sf.Frags = append(sf.Frags, common_peer.Frag{})
			copy(sf.Frags[i+1:], sf.Frags[i:])
			sf.Frags[i] = frag
			return true
		}
	}
	sf.Frags = append(sf.Frags, frag)
	return true
}

//Add does the same thing as AddFrag
func (sf *StoredFrags) Add(index uint64, start int64, size int64) bool {
	return sf.AddFrag(common_peer.Frag{index, start, size})
}

//Remove removes the given number
func (sf *StoredFrags) Remove(n uint64) bool {
	for i := 0; i < len(sf.Frags); i++ {
		if sf.Frags[i].Index < n {
			continue
		} else if sf.Frags[i].Index == n {
			if n == uint64(len(sf.Frags)-1) {
				//if n refers to the last element
				//then directly cut off the tail
				sf.Frags = sf.Frags[:len(sf.Frags)-1]
			} else {
				//n refers to one element of the middle elements
				copy(sf.Frags[i:], sf.Frags[i+1:])
				sf.Frags[len(sf.Frags)-1] = common_peer.Frag{}
				sf.Frags = sf.Frags[:len(sf.Frags)-1]
			}
			return true
		} else {
			return false
		}
	}
	return false
}

//Equal returns whether two StoredFrags are equal
func (sf *StoredFrags) Equal(sf2 *StoredFrags) bool {
	if len(sf.Frags) != len(sf2.Frags) {
		return false
	}
	for i, n := range sf.Frags {
		if sf2.Frags[i] != n {
			return false
		}
	}
	return true
}
