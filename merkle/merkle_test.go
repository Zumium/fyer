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

func TestMarshalUnmarshal(t *testing.T) {
	testFrag := [][]byte{
		[]byte("abc"),
		[]byte("def"),
		[]byte("ghi"),
		[]byte("jkl"),
		[]byte("mno"),
	}

	tree := NewTreeFrom(testFrag)
	b, err := Marshal(tree)
	if err != nil {
		t.Fatalf("marshaling failed: %s", err)
	}
	tree2, err := Unmarshal(b)
	if err != nil {
		t.Fatalf("unmarshaling failed: %s", err)
	}

	size1, size2 := tree.Size(), tree2.Size()
	if size1 != size2 {
		t.Fatalf("size not matching, %d-%d", size1, size2)
	}
	for i := uint64(0); i < size1; i++ {
		if !((tree.nodes[i] == nil && tree2.nodes[i] == nil) || bytes.Equal(tree.nodes[i], tree2.nodes[i])) {
			t.Fatalf("node not matching: %v -- %v", tree.nodes[i], tree2.nodes[i])
		}
	}
}
