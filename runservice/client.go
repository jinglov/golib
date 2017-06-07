package runservice

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

type runClient struct {
	mu      sync.Mutex
	net     string
	addr    string
	conn    net.Conn
	log     *log.Logger
	isClose bool
	isOpen  bool
}

var client *runClient

func NewClient(lnet, addr string) {
	var err error
	client = &runClient{
		net:  lnet,
		addr: addr,
		log:  log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
	client.conn, err = net.Dial(client.net, client.addr)
	if err != nil {
		log.Println(err)
		return
	}
}

func Send(cmd uint8, param string) error {
	return client.send(cmd, param)
}

func Close() {
	client.close()
}

func (c *runClient) close() {
	c.conn.Close()
}

func (c *runClient) send(cmd uint8, param string) error {
	//defer c.conn.Close()
	buf := bytes.NewBuffer(make([]byte, 0, len(param)+3))
	buf.WriteByte(cmd)
	length := len(param)
	binary.Write(buf, binary.BigEndian, uint16(length))
	buf.WriteString(param)
	c.log.Println(buf.Bytes())
	_, err := c.conn.Write(buf.Bytes())
	if err == nil {
		c.log.Println("response")
		resHead := make([]byte, 2)
		_, err = c.conn.Read(resHead)
		c.log.Println(resHead)
		var resLen uint16
		buf := bytes.NewBuffer(resHead)
		binary.Read(buf, binary.BigEndian, &resLen)
		if resLen > 0 {
			info := make([]byte, resLen)
			var rcount uint16
			for {
				/*
					读消息
				*/
				rlen, err := c.conn.Read(info)
				if err != nil && err == io.EOF {
					break
				}
				if err != nil {
					c.log.Println(err)
					return err
				}
				rcount += uint16(rlen)
				if rcount == resLen {
					break
				}
			}
			c.log.Println(info)
			c.log.Println(string(info))
		}
	}
	return err
}
