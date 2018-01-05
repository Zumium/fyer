package peer

import (
	"testing"

	"github.com/Zumium/fyer/common"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

func TestFileDBWrapper(t *testing.T) {
	//prepare db instance
	var err error
	dbInstance, err = leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		t.Fatal(err)
	}

	filename := "testfile"

	f := ToFile(filename)
	testStoredFrags1 := &StoredFrags{Frags: []common.Frag{
		{0, 0, 2047},
		{1, 2048, 4095},
	}}
	err = f.Edit().SetStoredFrags(testStoredFrags1).Done()
	if err != nil {
		t.Fatal(err)
	}

	f2 := ToFile(filename)
	testStoredFrags2 := f2.StoredFrags()
	if !testStoredFrags1.Equal(testStoredFrags2) {
		t.Fatal("not matching")
	}

	//tim := time.Now()
	//_, err = NewFileRecord(filename, 65530, []byte("abc"), 3, tim, mtree.NewTreeFrom([][]byte{[]byte("hello")}), &StoredFrags{Numbers: []uint64{1, 4, 9}})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//wrapper2 := NewFilesDBWrapper(filename)
	//
	//fc, err := wrapper2.FragCount()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if fc != 3 {
	//	t.Fatal("frag count not matching")
	//}
	//
	//h, err := wrapper2.Hash()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if !bytes.Equal(h, []byte("abc")) {
	//	t.Fatal("hash value not matching")
	//}
	//
	//mt, err := wrapper2.MerkleTree()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Logf("merkle tree size is %d\n", mt.Size())
	//if b1, b2 := sha3.Sum256([]byte("hello")), mt.CopyGet(0); !bytes.Equal(b2, b1[:]) {
	//	t.Fatalf("merkle tree not matching, is %v, should be %v\n", b2, b1)
	//}
	//
	//sz, err := wrapper2.Size()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if sz != 65530 {
	//	t.Fatal("size not matching")
	//}
	//
	//sf, err := wrapper2.StoredFrags()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if !(sf.Numbers[0] == 1 && sf.Numbers[1] == 4 && sf.Numbers[2] == 9) {
	//	t.Fatal("stored frags not matching")
	//}
	//
	//ti, err := wrapper2.UploadTime()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if !ti.Equal(tim) {
	//	t.Fatal("upload time not matching")
	//}
}
