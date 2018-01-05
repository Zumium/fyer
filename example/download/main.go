package main

import (
	control_fyerwork "github.com/Zumium/fyer/control/fyerwork"
	"fmt"
	"os"
	"github.com/Zumium/fyer/cfg"
)

func main() {
	fmt.Println("Initializing config")
	if err:=cfg.Init();err!=nil{
		fmt.Fprintf(os.Stderr, "Failed to initializing config: %s\n", err)
		os.Exit(-1)
	}
	if err:=control_fyerwork.Download("hello.txt", "/");err!=nil{
		fmt.Fprintf(os.Stderr, "Failed to download: %s\n", err)
		os.Exit(-1)
	}
}