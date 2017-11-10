package main

import (
	"fmt"
	"github.com/jinglov/golib/cmdservice"
)

func main() {
	RunService()
}

func RunService() {
	cmdservice.NewService("unix", "test.sock")
	cmdservice.ServerHandler(2, "hello", sayhello)
	cmdservice.Start()
	select {}
}

func sayhello(b []byte) []byte {
	fmt.Println(string(b))
	return b
}
