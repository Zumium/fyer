package alg

import (
	common_peer "github.com/Zumium/fyer/common/peer"
)

type FragCutter struct {
	Size     uint64
	FragSize int64
}

func (fragCutter *FragCutter) Cut() (frags []common_peer.Frag) {
	fragCount := fragCutter.Size / uint64(fragCutter.FragSize)
	frags = make([]common_peer.Frag, fragCount+1)
	var start int64
	for i := uint64(0); i < fragCount; i++ {
		frags[i].Index = i
		frags[i].Start = start
		frags[i].Size = fragCutter.FragSize
		start += fragCutter.FragSize
	}
	frags[fragCount].Index = fragCount
	frags[fragCount].Start = start
	frags[fragCount].Size = int64(fragCutter.Size) - start
	return
}
