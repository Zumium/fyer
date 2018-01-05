package alg

import (
	"github.com/Zumium/randutil"
	"math"
)

//DispatchAlg 文件段分派算法
//文件分派算法使用函数 p=-10^count 确定节点被分派文件段的比例
type DispatchAlg struct {
	fragCount uint64
	replica   int

	weights    []*randutil.Choice
	holdCounts map[string]uint32
}

//NewDispatchAlg creates a new algorithm state instance
func NewDispatchAlg(peers []string, fragCount uint64, replica int) *DispatchAlg {
	alg := &DispatchAlg{
		fragCount: fragCount,
		replica:   replica,

		holdCounts: make(map[string]uint32),
		weights:    make([]*randutil.Choice, 0, len(peers)),
	}
	for _, peer := range peers {
		alg.weights = append(alg.weights, &randutil.Choice{Weight: alg.weight(0), Item: peer})
		alg.holdCounts[peer] = 0
	}
	return alg
}

func (alg *DispatchAlg) weight(count uint32) int {
	//return int(1000 - 1000.0 / float64(alg.fragCount * uint64(alg.replica) / uint64(cap(alg.weights))) * float64(count))
	//return int(1000 * math.Exp(-float64(count))) + 10
	shouldHave := alg.fragCount * uint64(alg.replica) / uint64(cap(alg.weights))
	return int(math.Pow10(int(shouldHave)-int(count))) + 1
}

//selectGroup selects a group of different peers, the group size is equal to alg.replicate
func (alg *DispatchAlg) selectGroup() []string {
	group := make([]string, 0, alg.replica)

	var lastPeer string
	for j := 0; j < alg.replica; j++ {
		//TEST
		//fmt.Printf("%d %d %d\n", alg.weights[0].Weight, alg.weights[1].Weight, alg.weights[2].Weight)
		//fmt.Printf("%d %d %d\n", alg.holdCounts["A"], alg.holdCounts["B"], alg.holdCounts["C"])
		//TEST
		choice, err := randutil.WeightedChoice(alg.weights)
		if err != nil {
			return nil
		}
		if peer := choice.Item.(string); peer != lastPeer {
			group = append(group, peer)
			lastPeer = peer
			alg.holdCounts[peer]++
			choice.Weight = alg.weight(alg.holdCounts[peer])
		} else {
			j-- // selected same peer as last time, select one more time
		}
	}

	return group
}

//Dispatch run
//func (alg *DispatchAlg) Dispatch() map[uint64][]string {
//	result := make(map[uint64][]string)
//
//	for i := uint64(0); i < alg.fragCount; i++ {
//		result[i] = alg.selectGroup()
//	}
//
//	return result
//}
func (alg *DispatchAlg) Dispatch() [][]string {
	result := make([][]string, alg.fragCount)

	for i := uint64(0); i < alg.fragCount; i++ {
		result[i] = alg.selectGroup()
	}

	return result
}
