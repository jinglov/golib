package runservice

import "sync"

type runHandler struct {
	id      uint8
	name    string
	handler handlerFun
}

var initMu sync.Mutex

type handlerFun func(params []byte) []byte
type serverHandler map[uint8]*runHandler

func defaultServerHandler() serverHandler {
	res := make(serverHandler)
	res.Add(1, "ping", ping)
	return res
}

func (h serverHandler) Add(id uint8, name string, handler handlerFun) {
	initMu.Lock()
	defer initMu.Unlock()
	r := &runHandler{
		id:      id,
		name:    name,
		handler: handler,
	}
	h[id] = r
}

func ping(params []byte) []byte {
	return []byte("ok....")
	//return make([]byte, 0), nil
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
