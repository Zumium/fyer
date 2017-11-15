package fragmngr

import (
	"bytes"
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

	if err := fa.Write(2, testdata); err != nil {
		t.Fatal(err)
	}
	if err := fa.Write(4, testdata); err != nil {
		t.Fatal(err)
	}

	dout, err := fa.Read(2)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(dout, testdata) {
		t.Fatalf("fragment read not matching, is %v, should be %v\n", dout, testdata)
	}
	dout, err = fa.Read(4)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(dout, testdata) {
		t.Fatalf("fragment read not matching, is %v, should be %v\n", dout, testdata)
	}
}
