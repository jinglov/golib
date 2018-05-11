package cmdservice

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"log"
)

type cmdServer struct {
	mu      sync.Mutex
	net     string
	addr    string
	lis     net.Listener
	log     *log.Logger
	isClose bool
	isOpen  bool
	handler serverHandler
}

var server *cmdServer

func NewService(lnet, addr string) {
	server = &cmdServer{
		net:     lnet,
		addr:    addr,
		handler: defaultServerHandler(),
	}
}

func ServerHandler(id uint8, name string, handler HandlerFun) error {
	if server == nil {
		return errors.New("init service first.")
	}
	server.handler.Add(id, name, handler)
	return nil
}

func Start() error {
	if server != nil && !server.isClose {
		return server.start()
	}
	return nil
}

func (r *cmdServer) clean() {
	if r.net == "unix" {
		os.Remove(r.addr)
	}
}

func (r *cmdServer) close() (err error) {
	err = r.lis.Close()
	if err != nil {
		logError("Service stop error ", r.net, r.addr, err.Error())
		return
	}
	if r.net == "unix" {
		err = os.Remove(r.addr)
		if err != nil {
			logError("Remove sock file error:", err.Error())
			return
		}
	}
	info("Service stop at ", r.net, r.addr)
	return
}

func (r *cmdServer) start() (err error) {
	info("Service start at ", r.net, r.addr)
	r.clean()
	r.lis, err = net.Listen(r.net, r.addr)
	if err != nil {
		logError("Service start error:", err)
		return
	}
	go func() {
		server.isOpen = true
		defer func() {
			server.isOpen = false
			server.isClose = false
		}()
		defer r.close()
		for {
			if conn, err := r.lis.Accept(); err == nil {
				go newAccpet(conn)
			}
		}
	}()
	return
}

type accept struct {
	conn  net.Conn      //连接
	chcmd chan byte     //得到命令
	chr   chan []byte   //读数据缓存
	chw   chan []byte   //写数据缓存
	chend chan struct{} //结束标记
}

func newAccpet(conn net.Conn) {
	ac := &accept{
		conn:  conn,
		chcmd: make(chan byte),
		chr:   make(chan []byte),
		chw:   make(chan []byte),
		chend: make(chan struct{}),
	}
	ac.run()
}

func (c *accept) receive() {
	defer func() {
		info("close")
		c.chend <- struct{}{} //断开连接状态
	}()
	for {
		/*
			读消息头
		*/
		header := make([]byte, 3)
		_, err := c.conn.Read(header)
		debug(header)
		if err != nil && err == io.EOF {
			debug("EOF")
			return
		}
		//如果有错误，把错误抛出来，并且断开连接
		if err != nil {
			logError(err)
			return
		}
		cmd := header[0]
		debug(cmd)
		c.chcmd <- cmd
		/*
			读消息长度
		*/
		var length uint16
		buf := bytes.NewBuffer(header[1:3])
		binary.Read(buf, binary.BigEndian, &length)
		debug(length)
		/*
			读消息实体
		*/
		info := make([]byte, length)
		/*
			读消息
		*/
		_, err = c.conn.Read(info)
		//如果结束就返回
		if err != nil && err == io.EOF {
			debug("EOF")
			return
		}
		//如果有错误，把错误抛出来，并且断开连接
		if err != nil {
			logError(err)
			return
		}
		c.chr <- info
	}
	return
}

func (c *accept) send() {
	defer close(c.chw)
	buf := bytes.NewBuffer(make([]byte, 0))
	for {
		select {
		case response := <-c.chw:
			if response == nil { //写空结束
				return
			}
			info("response:", response)
			resLen := len(response)
			buf.Reset()
			buf.Grow(resLen + 1)
			binary.Write(buf, binary.BigEndian, uint16(resLen)) //写入消息长度
			buf.Write(response)                                 //写入消息内容
			if _, cerr := c.conn.Write(buf.Bytes()); cerr != nil {
				logError(cerr)
			}
		}
	}
}

func (c *accept) handler() {
	defer close(c.chr)
	defer close(c.chcmd)
	var response []byte
	for {
		select {
		case cmd := <-c.chcmd:
			debug(cmd)
			info := <-c.chr
			if info == nil { //接收到空
				return
			}
			/*
				执行命令
			*/
			if handler, ok := server.handler[cmd]; ok {
				response = handler.handler(info)
			} else {
				response = []byte("cmd:" + strconv.Itoa(int(cmd)) + " not support.")
			}
			if response == nil {
				response = make([]byte, 0)
			}
			c.chw <- response
		}
	}
}

func (c *accept) run() {
	defer close(c.chend)
	go c.receive()
	go c.send()
	go c.handler()
	select {
	case <-c.chend:
		c.chw <- nil //写空关闭handle 协程
		c.chr <- nil //读空关闭写协程
		c.chcmd <- 0 //命令0
		c.conn.Close()
		info("close conn")
		return
	}
}
