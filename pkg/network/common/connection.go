package common

import (
	"errors"
	"github.com/gorilla/websocket"
	"net"
	"sync"
	"time"
)

const (
	Tcp = iota
	Websocket
)

type Connection struct {
	types int
	Conn net.Conn
	TcpConn *net.TCPConn
	WebsocketConn *websocket.Conn

	ConnID uint32
	isClosed bool
	ExitBuffChan chan bool
	msgChan chan *Message
	msgBuffChan chan *Message
	maxPacketSize uint32

	property map[string]interface{}
	propertyLock sync.RWMutex
	addTime int64
	updateTime int64

	OnClose func(conn *Connection)
}

func NewConnection(types int, conn net.Conn, connID uint32, maxMsgChanLen uint32, maxPacketSize uint32) *Connection {
	c := &Connection{
		types:		  types,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		msgChan:  	  make(chan *Message),
		msgBuffChan:  make(chan *Message, maxMsgChanLen),
		maxPacketSize:maxPacketSize,
		property:     make(map[string]interface{}),
		addTime: 	  time.Now().Unix(),
	}

	return c
}

func (c *Connection) SetTcpConn(tcpConn *net.TCPConn) {
	c.TcpConn = tcpConn
}

func (c *Connection) SetWebsocketConn(websocketConn *websocket.Conn) {
	c.WebsocketConn = websocketConn
}

func (c *Connection) SetOnClose(hookFunc func(*Connection)) {
	c.OnClose = hookFunc
}

func (c *Connection) GetConnection() net.Conn {
	return c.Conn
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.TcpConn
}

func (c *Connection) GetWebsocketConnection() *websocket.Conn {
	return c.WebsocketConn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId byte, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}

	c.msgChan <- NewMsg(msgId, data)

	return nil
}

func (c *Connection) SendBuffMsg(msgId byte, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}

	c.msgBuffChan <- NewMsg(msgId, data)

	return nil
}

func (c *Connection) GetMsgChan()  <-chan *Message {
	return c.msgChan
}

func (c *Connection) GetMsgBuffChan() <-chan *Message {
	return c.msgBuffChan
}

func (c *Connection) Close() error {
	if c.isClosed == true {
		return nil
	}
	c.isClosed = true

	if c.OnClose != nil {
		c.OnClose(c)
	}

	err := c.Conn.Close()
	c.ExitBuffChan <- true

	close(c.ExitBuffChan)
	close(c.msgChan)
	close(c.msgBuffChan)

	return err
}

func (c *Connection) Closed() bool {
	return c.isClosed
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func (c *Connection) OnlineDuration() int64 {
	return c.updateTime - c.addTime
}

func (c *Connection) GetUpdateTime() int64 {
	return c.updateTime
}

func (c *Connection) UpdateTime() {
	c.updateTime = time.Now().Unix()
}