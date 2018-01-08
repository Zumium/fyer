package center

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"sort"
	"time"
)

func TestFileMetaModifying(t *testing.T) {
	viper.Set("mongo_address", "127.0.0.1:27017")

	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Close()

	testHash := []byte("hello world")

	//create record
	fmeta, _ := ToFileMeta("test")
	if err := fmeta.Edit().SetFragCount(3).SetHash(testHash).SetSize(254).SetUploadTime(time.Now()).Done(); err != nil {
		t.Fatal(err)
	}
	defer fmeta.Remove()

	//check storage success
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

	//check modifying success
	fmeta3, _ := ToFileMeta("test")
	if err := fmeta3.Edit().SetFragCount(5).SetSize(77).Done(); err != nil {
		t.Fatal(err)
	}
	fmeta4, _ := ToFileMeta("test")
	fragCount4 := fmeta4.FragCount()
	if fragCount4 != 5 {
		t.Fatalf("frag count mismatching: is %d, should be %d\n", fragCount, 3)
	}
	hash4 := fmeta4.Hash()
	if !bytes.Equal(hash4, testHash) {
		t.Fatalf("hash mismatching: is %v, should be %v\n", hash, testHash)
	}
	size4 := fmeta4.Size()
	if size4 != 77 {
		t.Fatalf("size mismatching: is %d, should be %d\n", size, 254)
	}
	uploadTime4 := fmeta4.UploadTime()
	if uploadTime4.IsZero() {
		t.Fatal("upload time is zero")
	}
}

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

func TestGetFileNames(t *testing.T) {
	viper.Set("mongo_address", "127.0.0.1:27017")

	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Close()

	testNames := []string{"A", "B", "C"}

	for _, name := range testNames {
		handle, _ := ToFileMeta(name)
		if err := handle.Edit().Done(); err != nil {
			t.Fatal(err)
		}
	}

	names, err := Files()
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(names)
	if len(names) != len(testNames) {
		t.Fatal("files count not matching")
	}
	for i, n := range testNames {
		if names[i] != n {
			t.Fatalf("file name not correct at index %d: is %s, should be %s", i, names[i], n)
		}
	}
}
