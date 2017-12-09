package alg

import (
	"testing"
)

//import "testing"
//
//func TestDispatchAlg(t *testing.T) {
//	l := []string{"A", "B", "C", "D"}
//	fragCount := uint64(20)
//
//	alg := NewDispatchAlg(l, fragCount)
//	t.Log(alg.Dispatch())
//	t.Log(alg.Dispatch())
//	t.Log(alg.Dispatch())
//	t.Log(alg.Dispatch())
//	t.Fail()
//}

func TestDispatchAlg(t *testing.T) {
	peers := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	alg := NewDispatchAlg(peers, 10, 2)
	dispatchTable := alg.Dispatch()

	t.Log(dispatchTable) //目视检查

	statics := make(map[string]int)
	for _, peer := range peers {
		statics[peer] = 0
	}
	for _, peers := range dispatchTable {
		for _, peer := range peers {
			statics[peer]++
		}
	}
	t.Fatal(statics) //目视检查
}
