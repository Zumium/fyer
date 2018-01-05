package peer

import (
	common "github.com/Zumium/fyer/common"
)

//StoredFrags represents fragment index numbers which has been stored locally
// type StoredFrags []uint64
type StoredFrags struct {
	Frags []common.Frag `json:"frags"`
}

//NewEmptyStoredFrags creates an new empty StoredFrags struct
func NewEmptyStoredFrags() *StoredFrags {
	return new(StoredFrags)
}

func (sf *StoredFrags) binarySearch(n uint64) int {
	l := len(sf.Frags)
	switch l {
	case 0:
		return -1
	case 1:
		if sf.Frags[0].Index == n {
			return 0
		} else {
			return -1
		}
	}
	if n < sf.Frags[0].Index || n > sf.Frags[l-1].Index {
		return -1
	}
	//二分查找
	var (
		begin int
		end   = l - 1
		mid   int
	)
	for begin < end {
		if sf.Frags[begin].Index == n {
			return begin
		}
		if sf.Frags[end].Index == n {
			return end
		}
		mid = (begin + end) / 2
		if n < sf.Frags[mid].Index {
			end = mid
		} else {
			begin = mid + 1
		}
	}
	return -1
}

//Has returns true if the n-th fragment has been stored locally
func (sf *StoredFrags) Has(n uint64) bool {
	return sf.binarySearch(n) != -1
}

func (sf *StoredFrags) Find(n uint64) (common.Frag, bool) {
	pos := sf.binarySearch(n)
	if pos == -1 {
		return common.Frag{}, false
	}
	return sf.Frags[pos], true
}

//AddFrag inserts a new Frag while keep it sorted, returns the result whether
//the frag has been successfully inserted
func (sf *StoredFrags) AddFrag(frag common.Frag) bool {
	for i := 0; i < len(sf.Frags); i++ {
		if sf.Frags[i].Index < frag.Index {
			continue
		} else if sf.Frags[i].Index == frag.Index {
			return false
		} else {
			sf.Frags = append(sf.Frags, common.Frag{})
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
	return sf.AddFrag(common.Frag{index, start, size})
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
				sf.Frags[len(sf.Frags)-1] = common.Frag{}
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
