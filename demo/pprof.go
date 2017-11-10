package main

import "github.com/jinglov/golib/cmdpprof"

func main() {
	//demoServer()
	demoClient()
}

func demoServer() {
	cmdpprof.NewPprofServer("unix", "test.sock")
	select {}
}

func demoClient() {
	cmdpprof.NewCmdClient("unix", "test.sock")
}
