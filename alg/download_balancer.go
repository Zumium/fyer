package alg

import (
	"container/list"
)

type peerRank struct {
	Peer  string
	Count int
}

type DownloadBalancer struct {
	fragDistribution [][]string

	peers map[string]bool
	rank  *list.List
}

func NewDownloadBalancer(fragDistribution [][]string) *DownloadBalancer {
	balancer := &DownloadBalancer{fragDistribution: fragDistribution, peers: make(map[string]bool), rank: list.New()}
	return balancer
}

func (balancer *DownloadBalancer) Result() []string {
	result := make([]string, len(balancer.fragDistribution))
	for i, peers := range balancer.fragDistribution {
		//add non-exist peer string into rank
		for _, peer := range peers {
			if !balancer.peers[peer] {
				balancer.rank.PushFront(&peerRank{Peer: peer, Count: 0})
				balancer.peers[peer] = true
			}
		}
		//fetch data from first found peer
		var selectedElem *list.Element
	SelectPeer:
		for e := balancer.rank.Front(); e != nil; e = e.Next() {
			pr := e.Value.(*peerRank)
			for _, peer := range peers {
				if pr.Peer == peer {
					result[i] = peer
					pr.Count++
					selectedElem = e
					break SelectPeer
				}
			}
		}

		//increase the found peer's count, then re-rank by count increasing
		if selectedElem != nil {
			targetElem := selectedElem.Next()
			for ; targetElem != nil && targetElem.Value.(*peerRank).Count < selectedElem.Value.(*peerRank).Count; targetElem = targetElem.Next() {
			}
			if targetElem == nil {
				//move to the end
				balancer.rank.MoveToBack(selectedElem)
			} else {
				balancer.rank.MoveBefore(selectedElem, targetElem)
			}
		}
	}

	return result
}
