package cmdservice

import "sync"

type cmdHandler struct {
	id      uint8
	name    string
	handler HandlerFun
}

var initMu sync.Mutex

type HandlerFun func(params []byte) []byte
type serverHandler map[uint8]*cmdHandler

func defaultServerHandler() serverHandler {
	res := make(serverHandler)
	res.Add(1, "ping", ping)
	return res
}

func (h serverHandler) Add(id uint8, name string, handler HandlerFun) {
	initMu.Lock()
	defer initMu.Unlock()
	r := &cmdHandler{
		id:      id,
		name:    name,
		handler: handler,
	}
	h[id] = r
}

func ping(params []byte) []byte {
	return []byte("pong")
}

type clientHandler map[string]uint8

func defaultClientHandler() clientHandler {
	res := make(clientHandler)
	res.Add(1, "ping")
	return res
}

func (h clientHandler) Add(id uint8, name string) {
	initMu.Lock()
	defer initMu.Unlock()
	h[name] = id
}
