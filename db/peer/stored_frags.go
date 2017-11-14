package peer

//StoredFrags represents fragment index numbers which has been stored locally
// type StoredFrags []uint64
type StoredFrags struct {
	Numbers []uint64 `json:"numbers"`
}

//NewEmptyStoredFrags creates an new empty StoredFrags struct
func NewEmptyStoredFrags() *StoredFrags {
	return new(StoredFrags)
}

//Has returns true if the n-th fragment has been stored locally
func (sf *StoredFrags) Has(n uint64) bool {
	l := len(sf.Numbers)
	switch l {
	case 0:
		return false
	case 1:
		return sf.Numbers[0] == n
	}
	if n < sf.Numbers[0] || n > sf.Numbers[l-1] {
		return false
	}
	//二分查找
	var (
		begin uint64
		end   = uint64(l - 1)
		mid   uint64
	)
	for begin < end {
		if sf.Numbers[begin] == n || sf.Numbers[end] == n {
			return true
		}
		mid = (begin + end) / 2
		if n < sf.Numbers[mid] {
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
	for i := 0; i < len(sf.Numbers); i++ {
		if sf.Numbers[i] < n {
			continue
		} else if sf.Numbers[i] == n {
			return false
		} else {
			sf.Numbers = append(sf.Numbers, n)
			copy(sf.Numbers[i+1:], sf.Numbers[i:])
			sf.Numbers[i] = n
			return true
		}
	}
	sf.Numbers = append(sf.Numbers, n)
	return true
}

//Remove removes the given number
func (sf *StoredFrags) Remove(n uint64) bool {
	for i := 0; i < len(sf.Numbers); i++ {
		if sf.Numbers[i] < n {
			continue
		} else if sf.Numbers[i] == n {
			if n == uint64(len(sf.Numbers)-1) {
				sf.Numbers = sf.Numbers[:len(sf.Numbers)-1]
			} else {
				copy(sf.Numbers[i:], sf.Numbers[i+1:])
				sf.Numbers[len(sf.Numbers)-1] = 0
				sf.Numbers = sf.Numbers[:len(sf.Numbers)-1]
			}
			return true
		} else {
			return false
		}
	}
	return false
}
