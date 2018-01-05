package filemanager

import (
	"bytes"
	"os"
	"testing"
)

func TestFileManager(t *testing.T) {
	f, err := os.Create("/tmp/hello.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("Hello World!!!!!!!!!!!!!!!!!!!!!!!!")
	f.Close()

	if err := Register("/tmp/hello.txt"); err != nil {
		t.Fatal(err)
	}
	for name, path := range files {
		if name != "hello.txt" || path != "/tmp/hello.txt" {
			t.Fatalf("not correctly been stored, actual name is %s, path is %s", name, path)
		}
	}

	fi, err := Open("hello.txt")
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 50)
	n, err := fi.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf[:n], []byte("Hello World!!!!!!!!!!!!!!!!!!!!!!!!")) {
		t.Fatalf("wrong file content, actual data is %v", buf[:n])
	}
}
