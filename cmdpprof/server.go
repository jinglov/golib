package cmdpprof

import (
	"bytes"
	"github.com/jinglov/golib/cmdservice"
	"log"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"
)

func NewPprofServer(net, addr string) (err error) {
	cmdservice.NewService(net, addr)
	for _, p := range pprofHandlers {
		cmdservice.ServerHandler(p.id, p.name, p.cmdHandler)
	}
	return cmdservice.Start()
}

func cpuProfile(params []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	pprof.StartCPUProfile(buf)
	log.Println("start cpu pprof...")
	sec, _ := strconv.Atoi(string(params))
	if sec <= 0 {
		sec = 60
	}
	log.Println("sleep: ", sec)
	time.Sleep(time.Second * time.Duration(sec))
	pprof.StopCPUProfile()
	log.Println("stop cpu pprof...")
	return buf.Bytes()
}

func memProfile(param []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	log.Println("start memory pprof...")
	pprof.WriteHeapProfile(buf)
	log.Println("end memory pprof...")
	return buf.Bytes()
}

func blockProfile(params []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	log.Println("start block pprof...")
	pprof.Lookup("block").WriteTo(buf, 0)
	log.Println("stop block pprof...")
	return buf.Bytes()
}

func goroutineProfile(params []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	log.Println("start goroutine pprof...")
	pprof.Lookup("goroutine").WriteTo(buf, 0)
	log.Println("end goroutine pprof...")
	return buf.Bytes()
}

func threadProfile(params []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	log.Println("start thread pprof...")
	pprof.Lookup("thread").WriteTo(buf, 0)
	log.Println("end thread pprof...")
	return buf.Bytes()
}

func heapProfile(params []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	log.Println("start heap pprof...")
	pprof.Lookup("heap").WriteTo(buf, 0)
	log.Println("end heap pprof...")
	return buf.Bytes()
}
