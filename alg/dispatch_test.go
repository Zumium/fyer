package alg

import "testing"

func TestDispatchAlg(t *testing.T) {
	l := []string{"A", "B", "C", "D"}
	fragCount := uint64(20)

	alg := NewDispatchAlg(l, fragCount)
	t.Log(alg.Dispatch())
	t.Log(alg.Dispatch())
	t.Log(alg.Dispatch())
	t.Log(alg.Dispatch())
	t.Fail()
}
