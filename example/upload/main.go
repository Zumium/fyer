package main

import (
	"fmt"
	"github.com/Zumium/fyer/cfg"
	control_fyerwork "github.com/Zumium/fyer/control/fyerwork"
	"github.com/Zumium/fyer/filemanager"
	rpc_fyerwork "github.com/Zumium/fyer/rpc/fyerwork"
	"os"
)

func createHelloTxt() error {
	f, err := os.Create("/tmp/hello.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString("Hello World!!!!!!!!!!!!!"); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := createHelloTxt(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create /tmp/hello.txt: %s", err)
		os.Exit(-1)
	}

	fmt.Println("Initializing config")
	if err := cfg.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error happened: %s\n", err)
		os.Exit(-1)
	}

	fmt.Println("Registering test file -- /tmp/hello.txt")
	filemanager.Register("/tmp/hello.txt")
	go func() {
		if err := rpc_fyerwork.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error happened: %s\n", err)
			os.Exit(-1)
		}
	}()
	control_fyerwork.UploadFile("hello.txt", 24, []byte("abc"), "172.18.0.8")
	rpc_fyerwork.Wait()
}
