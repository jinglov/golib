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
		description:   "取CPU分析",
		cmdHandler:    cpuProfile,
		clientHandler: cpuClient,
	},
	{
		id:            11,
		name:          "mem",
		description:   "取内存分析",
		cmdHandler:    memProfile,
		clientHandler: memClient,
	},
}
