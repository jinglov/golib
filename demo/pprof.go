package main

import (
	"github.com/jinglov/golib/cmdpprof"
	"flag"
	"fmt"
	"os/signal"
	"os"
	"syscall"
)

func main() {
	action := flag.String("action", "client", "-action server[client]")
	socket := flag.String("socket", "test.sock", "-socket test.sock")
	flag.Parse()
	switch *action {
	case "server":
		demoServer(*socket)
	case "client":
		demoClient(*socket)
	default:
		fmt.Println("-action server or -action client")
	}
}

func demoServer(socket string) {
	cmdpprof.NewPprofServer("unix", socket)
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGINT)
	<-sign
	os.Exit(0)
}

func demoClient(socket string) {
	cmdpprof.NewCmdClient("unix", socket)
}
