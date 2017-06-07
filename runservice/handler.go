package runservice

import "sync"

type runHandler struct {
	id      uint8
	name    string
	handler func(params []byte) ([]byte, error)
}

var idHandlerMap = make(map[uint8]*runHandler)
var nameHandlerMap = make(map[string]*runHandler)
var initMu sync.Mutex

func AddHandler(id uint8, name string, handler func(params []byte) ([]byte, error)) {
	initMu.Lock()
	defer initMu.Unlock()
	r := &runHandler{
		id:      id,
		name:    name,
		handler: handler,
	}
	idHandlerMap[id] = r
	if name != "" {
		nameHandlerMap[name] = r
	}
}

func defaultHandler() {
	AddHandler(0, "stop", Stop)
	AddHandler(1, "ping", ping)
}

func Stop(params []byte) ([]byte, error) {
	runner.isClose = true
	err := runner.close()
	return nil, err
}

func ping(params []byte) ([]byte, error) {
	return []byte("ok...."), nil
	//return make([]byte, 0), nil
}
