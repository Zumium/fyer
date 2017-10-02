package merkle

import (
	"bytes"
	"testing"

	"golang.org/x/crypto/sha3"
)

func TestMerkleTreeBuildAndSelfValid(t *testing.T) {
	testFrag := [][]byte{
		[]byte("abc"),
		[]byte("def"),
		[]byte("ghi"),
		[]byte("jkl"),
		[]byte("mno"),
	}
	merkleTree := NewTreeFrom(testFrag)
	if !merkleTree.SelfValid() {
		t.Fatal("the tree build is invalid")
	}
}

func TestLeavesReturnsCorrectly(t *testing.T) {
	testFrag := [][]byte{
		[]byte("abc"),
		[]byte("def"),
		[]byte("ghi"),
		[]byte("jkl"),
		[]byte("mno"),
	}
	tree := NewTreeFrom(testFrag)
	leaves := tree.Leaves()
	if l1, l2 := len(testFrag), len(leaves); l1 != l2 {
		t.Fatalf("the merkle tree has incorrect number of leaves: %d - %d", l1, l2)
	}
	for i := 0; i < len(testFrag); i++ {
		if correctHash := sha3.Sum256(testFrag[i]); !bytes.Equal(correctHash[:], leaves[i]) {
			t.Fatalf("leave not matching at index %d: is %v, should be %v", i, leaves[i], correctHash)
		}
	}
}
