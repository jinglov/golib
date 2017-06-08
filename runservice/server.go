package runservice

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

type runService struct {
	mu      sync.Mutex
	net     string
	addr    string
	lis     net.Listener
	log     *log.Logger
	isClose bool
	isOpen  bool
	handler serverHandler
}

var runner *runService

func NewService(lnet, addr string) {
	runner = &runService{
		net:     lnet,
		addr:    addr,
		log:     log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
		handler: defaultServerHandler(),
	}
}

func ServerHandler(id uint8, name string, handler handlerFun) error {
	if runner == nil {
		return errors.New("init service first.")
	}
	runner.handler.Add(id, name, handler)
	return nil
}

func Start() {
	if runner != nil && !runner.isClose {
		go runner.start()
	}
}

func (r *runService) clean() {
	if r.net == "unix" {
		os.Remove(r.addr)
	}
}

func (r *runService) close() (err error) {
	err = r.lis.Close()
	if err != nil {
		r.log.Println("Service stop error ", r.net, r.addr, err.Error())
		return
	}
	if r.net == "unix" {
		err = os.Remove(r.addr)
		if err != nil {
			r.log.Println("Remove sock file error:", err.Error())
			return
		}
	}
	r.log.Println("Service stop at ", r.net, r.addr)
	return
}

func (r *runService) start() {
	var err error
	r.log.Println("Service start at ", r.net, r.addr)
	r.clean()
	r.lis, err = net.Listen(r.net, r.addr)
	if err != nil {
		r.log.Println("Service start error:", err)
		return
	}
	runner.isOpen = true
	defer func() {
		runner.isOpen = false
		runner.isClose = false
	}()
	defer r.close()
	for {
		if conn, err := r.lis.Accept(); err == nil {
			go newAccpet(conn).run()
		}
	}
}

type accept struct {
	conn  net.Conn
	chcmd chan byte
	chr   chan []byte
	chw   chan []byte
	chend chan struct{}
}

func newAccpet(conn net.Conn) *accept {
	return &accept{
		conn:  conn,
		chcmd: make(chan byte),
		chr:   make(chan []byte),
		chw:   make(chan []byte),
		chend: make(chan struct{}),
	}
}

func (c *accept) receive() {
	defer func() {
		runner.log.Println("close")
		c.chend <- struct{}{} //断开连接状态
	}()
	for {
		/*
			读消息头
		*/
		header := make([]byte, 3)
		_, err := c.conn.Read(header)
		runner.log.Println(header)
		if err != nil && err == io.EOF {
			runner.log.Println("EOF")
			return
		}
		//如果有错误，把错误抛出来，并且断开连接
		if err != nil {
			runner.log.Println(err)
			return
		}
		cmd := header[0]
		runner.log.Println(cmd)
		c.chcmd <- cmd
		/*
			读消息长度
		*/
		var length uint16
		buf := bytes.NewBuffer(header[1:3])
		binary.Read(buf, binary.BigEndian, &length)
		//runner.log.Println(length)
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
			runner.log.Println("EOF")
			return
		}
		//如果有错误，把错误抛出来，并且断开连接
		if err != nil {
			runner.log.Println(err)
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
			runner.log.Println("response:", response)
			resLen := len(response)
			buf.Reset()
			buf.Grow(resLen + 1)
			binary.Write(buf, binary.BigEndian, uint16(resLen)) //写入消息长度
			buf.Write(response)                                 //写入消息内容
			if _, cerr := c.conn.Write(buf.Bytes()); cerr != nil {
				runner.log.Println(cerr)
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
			runner.log.Println(cmd)
			info := <-c.chr
			if info == nil { //接收到空
				return
			}
			/*
				执行命令
			*/
			if handler, ok := runner.handler[cmd]; ok {
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
		runner.log.Println("close conn")
		return
	}
}
