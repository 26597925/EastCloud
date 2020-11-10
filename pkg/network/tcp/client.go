package tcp

import (
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/26597925/EastCloud/pkg/network/common"
	"io"
	"net"
)

type Client struct {
	conn *common.Connection
	msgHandler *common.MsgHandle
}

func NewClient() *Client {
	s := &Client{
		msgHandler: common.NewMsgHandle(10, 1024),
	}
	return s
}

func (c *Client) Connect() {
	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		logger.ErrorF("common failed, err : %v\n", err.Error())
		return
	}
	defer conn.Close()

	c.conn = common.NewConnection(common.Tcp, conn, 0, 1024, 4096)

	c.msgHandler.StartWorkerPool()
	go c.startReader()
    c.startWriter()
}

func (c *Client) GetConnect() *common.Connection {
	return c.conn
}

func (c *Client) SendMsg(msgId byte, data []byte) error {
	return c.conn.SendMsg(msgId, data)
}

func (c *Client) SendBuffMsg(msgId byte, data []byte) error {
	return c.conn.SendBuffMsg(msgId, data)
}

func (c *Client) AddRouter(msgId byte, router common.IRouter) {
	c.msgHandler.AddRouter(msgId, router)
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) startReader() {
	logger.Info("[Reader Goroutine is running]")
	defer logger.Info(c.conn.RemoteAddr().String(), "[conn Reader exit!]")
	defer c.Close()

	for {
		//c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		dp := newDataPack(4096)

		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.conn.Conn, headData); err != nil {
			logger.Error("client read msg head error ", err)
			break
		}

		msg, err := dp.Unpack(headData)
		if err != nil {
			logger.Error("unpack error ", err)
			break
		}

		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.conn.Conn, data); err != nil {
				logger.Error("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)

		req := common.Request{
			Conn: c.conn,
			Msg:  msg,
		}

		go c.msgHandler.DoReceive(&req)
	}
}

func (c *Client) startWriter() {
	logger.Info("[Writer Goroutine is running]")
	defer logger.Info(c.conn.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case msg, ok := <-c.conn.GetMsgChan():
			if ok {
				dp := newDataPack(4096)
				data, err := dp.Pack(msg)
				if err != nil {
					logger.Error("Pack error msg id = ", msg.Id)
					return
				}

				go func() {
					if _, err := c.conn.Conn.Write(data); err != nil {
						c.Close()
						logger.Error("Send Data error:, ", err, " Conn Writer exit")
						return
					}
				}()

			}
			//fmt.Printf("Send data succ! data = %+v\n", data)
		case msg, ok := <-c.conn.GetMsgBuffChan():
			if ok {
				dp := newDataPack(4096)
				data, err := dp.Pack(msg)
				if err != nil {
					logger.Error("Pack error msg id = ", msg.Id)
					return
				}
				if _, err := c.conn.Conn.Write(data); err != nil {
					c.Close()
					logger.Error("Client Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				logger.Info("msgBuffChan is Closed")
				break
			}
		case <-c.conn.ExitBuffChan:
			return
		}
	}
}