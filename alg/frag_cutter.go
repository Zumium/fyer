package alg

import (
	common "github.com/Zumium/fyer/common"
)

type FragCutter struct {
	Size     uint64
	FragSize int64
}

func NewFragCutter(size uint64, fragSize int64) *FragCutter {
	return &FragCutter{size, fragSize}
}

func (fragCutter *FragCutter) Cut() (frags []common.Frag) {
	fragCount := fragCutter.Size / uint64(fragCutter.FragSize)
	frags = make([]common.Frag, fragCount+1)
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
