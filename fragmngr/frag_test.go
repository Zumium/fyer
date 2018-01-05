package fragmngr

import (
	"bytes"
	"github.com/Zumium/fyer/common"
	"os"
	"testing"
)

func TestFragMngr(t *testing.T) {
	testBaseDir := "/tmp/zumium/fyer/fragmngr/base"
	if err := os.MkdirAll(testBaseDir, 0777); err != nil {
		t.Fatal(err)
	}

	if err := InitSimpleFSFragManager(testBaseDir); err != nil {
		t.Fatal(err)
	}

	fa, err := FMInstance().Open("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer fa.Close()

	testdata := []byte("hello world")
	testdataSize := int64(len(testdata))

	testFragIndex0 := common.Frag{0, 0, testdataSize}
	testFragIndex2 := common.Frag{2, 2*testdataSize - 1, testdataSize}

	if err := fa.Write(testFragIndex0, testdata); err != nil {
		t.Fatal(err)
	}
	if err := fa.Write(testFragIndex2, testdata); err != nil {
		t.Fatal(err)
	}

	dout, err := fa.Read(testFragIndex0)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(dout, testdata) {
		t.Fatalf("fragment read not matching, is %v, should be %v\n", dout, testdata)
	}
	dout, err = fa.Read(testFragIndex2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(dout, testdata) {
		t.Fatalf("fragment read not matching, is %v, should be %v\n", dout, testdata)
	}
}
