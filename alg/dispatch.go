package alg

import (
	"math/rand"
)

//DispatchAlg 文件段分派算法
type DispatchAlg struct {
	fragCount uint64

	Labels    []string
	holdCount []uint32
}

//NewDispatchAlg 新建一个分派算法实例
func NewDispatchAlg(labels []string, fragCount uint64) *DispatchAlg {
	return &DispatchAlg{
		fragCount: fragCount,

		Labels:    labels,
		holdCount: make([]uint32, len(labels)),
	}
}

//Dispatch returns a label and assume the peer holding that label will
//hold a fragment
func (alg *DispatchAlg) Dispatch() (int, string) {
	peerCount := len(alg.Labels)
	// p:=make([]float64, peerCount)

	randNum := rand.Int31n(100)
	sum := 0

	idx := peerCount - 1
	for i := 0; i < peerCount; i++ {
		// p[i]=float64(1)/peerCount-float64(alg.holdCount[i])/alg.fragCount
		p := float64(1)/float64(peerCount) - float64(alg.holdCount[i])/float64(alg.fragCount)
		//Protect the algorithm when probability < 0
		if p < 0 {
			p = 0
		}
		if sum += int(100 * p); int32(sum) > randNum {
			idx = i
			break
		}
	}

	return idx, alg.Labels[idx]
}
