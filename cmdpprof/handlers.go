package cmdpprof

import "github.com/jinglov/golib/cmdservice"

type pprofHandler struct {
	id            uint8
	name          string
	description   string
	cmdHandler    cmdservice.HandlerFun
	clientHandler clientHandler
}

var pprofHandlers = []*pprofHandler{
	{
		id:            10,
		name:          "cpu",
		description:   "CPU分析",
		cmdHandler:    cpuProfile,
		clientHandler: cpuClient,
	},
	{
		id:            11,
		name:          "mem",
		description:   "内存分析",
		cmdHandler:    memProfile,
		clientHandler: memClient,
	},
	{
		id:            12,
		name:          "block",
		description:   "阻塞分析",
		cmdHandler:    blockProfile,
		clientHandler: blockClient,
	},
	{
		id:            13,
		name:          "goroutine",
		description:   "Goroutine分析",
		cmdHandler:    goroutineProfile,
		clientHandler: goroutineClient,
	},
	{
		id:            14,
		name:          "threadcreate",
		description:   "进程分析",
		cmdHandler:    threadProfile,
		clientHandler: threadClient,
	},
	{
		id:            15,
		name:          "heap",
		description:   "堆栈分析",
		cmdHandler:    heapProfile,
		clientHandler: heapClient,
	},
}
