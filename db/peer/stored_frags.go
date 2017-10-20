package peer

import (
	"encoding/json"
)

//StoredFrags represents fragment index numbers which has been stored locally
// type StoredFrags []uint64
type StoredFrags struct {
	numbers []uint64
}

//NewEmptyStoredFrags creates an new empty StoredFrags struct
func NewEmptyStoredFrags() *StoredFrags {
	return new(StoredFrags)
}

//Has returns true if the n-th fragment has been stored locally
func (sf *StoredFrags) Has(n uint64) bool {
	l := len(sf.numbers)
	switch l {
	case 0:
		return false
	case 1:
		return sf.numbers[0] == n
	}
	if n < sf.numbers[0] || n > sf.numbers[l-1] {
		return false
	}
	//二分查找
	var (
		begin uint64
		end   uint64 = l - 1
		mid   uint64
	)
	for begin < end {
		if sf.numbers[begin] == n || sf.numbers[end] == n {
			return true
		}
		mid = (begin + end) / 2
		if n < sf.numbers[mid] {
			end = mid
		} else {
			begin = mid + 1
		}
	}
	return false
}

//Add inserts a new number while keep it sorted, returns the result whether
//the number has been successfully inserted
func (sf *StoredFrags) Add(n uint64) bool {
	for i := 0; i < len(sf.numbers); i++ {
		if sf.numbers[i] < n {
			continue
		} else if sf.numbers[i] == n {
			return false
		} else {
			sf.numbers = append(sf.numbers, n)
			copy(sf.numbers[i+1:], sf.numbers[i:])
			sf.numbers[i] = n
			return true
		}
	}
	sf.numbers = append(sf.numbers, n)
	return true
}

//Remove removes the given number
func (sf *StoredFrags) Remove(n uint64) bool {
	for i := 0; i < len(sf.numbers); i++ {
		if sf.numbers[i] < n {
			continue
		} else if sf.numbers[i] == n {
			if n == len(sf.numbers)-1 {
				sf.numbers = sf.numbers[:len(sf.numbers)-1]
			} else {
				copy(sf.numbers[i:], sf.numbers[i+1:])
				sf.numbers[len(sf.numbers)-1] = nil
				sf.numbers = sf.numbers[:len(sf.numbers)-1]
			}
			return true
		} else {
			return false
		}
	}
	return false
}

//MarshalJSON implements json.Marshaler interface
func (sf *StoredFrags) MarshalJSON() ([]byte, error) {
	return json.Marshal(sf.numbers)
}

//UnmarshalJSON implements json.Unmarshaler interface
func (sf *StoredFrags) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &sf.numbers)
}
