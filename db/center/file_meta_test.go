package center

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
)

func TestFileMetaOperations(t *testing.T) {
	viper.Set("mongo_address", "127.0.0.1:27017")

	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Close()

	testHash := []byte("hello world")

	fmeta := ToFileMeta("test")
	if err := fmeta.Edit().SetFragCount(3).SetHash(testHash).SetSize(254).Done(); err != nil {
		t.Fatal(err)
	}
	defer fmeta.Remove()

	fmeta2 := ToFileMeta("test")
	fragCount, err := fmeta2.FragCount()
	if err != nil {
		t.Fatal(err)
	}
	if fragCount != 3 {
		t.Fatalf("frag count mismatching: is %d, should be %d\n", fragCount, 3)
	}
	hash, err := fmeta2.Hash()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(hash, testHash) {
		t.Fatalf("hash mismatching: is %v, should be %v\n", hash, testHash)
	}
	size, err := fmeta2.Size()
	if err != nil {
		t.Fatal(err)
	}
	if size != 254 {
		t.Fatalf("size mismatching: is %d, should be %d\n", size, 254)
	}
}
