package center

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"time"
)

func TestFileMetaOperations(t *testing.T) {
	viper.Set("mongo_address", "127.0.0.1:27017")

	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Close()

	testHash := []byte("hello world")

	fmeta, _ := ToFileMeta("test")
	if err := fmeta.Edit().SetFragCount(3).SetHash(testHash).SetSize(254).SetUploadTime(time.Now()).Done(); err != nil {
		t.Fatal(err)
	}
	defer fmeta.Remove()

	fmeta2, _ := ToFileMeta("test")
	fragCount := fmeta2.FragCount()
	if fragCount != 3 {
		t.Fatalf("frag count mismatching: is %d, should be %d\n", fragCount, 3)
	}
	hash := fmeta2.Hash()
	if !bytes.Equal(hash, testHash) {
		t.Fatalf("hash mismatching: is %v, should be %v\n", hash, testHash)
	}
	size := fmeta2.Size()
	if size != 254 {
		t.Fatalf("size mismatching: is %d, should be %d\n", size, 254)
	}
	uploadTime := fmeta2.UploadTime()
	if uploadTime.IsZero() {
		t.Fatal("upload time is zero")
	}
}
