package alg

import (
	"sort"
	"testing"
)

func TestDownloadBalancer(t *testing.T) {
	downloadBalancer := NewDownloadBalancer([][]string{
		[]string{"A", "B", "C"},
		[]string{"C", "A", "B"},
		[]string{"C", "B", "A"},
	})
	r := downloadBalancer.Result()
	for e := downloadBalancer.rank.Front(); e != nil; e = e.Next() {
		pr := e.Value.(*peerRank)
		if pr.Count != 1 {
			t.Fatalf("peerRank for peer %s is not 1", pr.Peer)
		}
	}
	sort.Strings(r)
	if r[0] != "A" || r[1] != "B" || r[2] != "C" {
		t.Fatal("result not corrent")
	}
}
