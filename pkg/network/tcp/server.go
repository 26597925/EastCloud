package tcp

import (
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/26597925/EastCloud/pkg/network/common"
	"github.com/26597925/EastCloud/pkg/util/atomic"
	"github.com/26597925/EastCloud/pkg/util/timer"
	"io"
	"net"
	"time"
)

type Server struct {
	msgHandler *common.MsgHandle
	connMgr *common.ConnManager
	tw *timer.TimingWheel

	OnConnStart func(conn *common.Connection)
	OnConnStop func(conn *common.Connection)
}

func NewServer() *Server {
	tw := timer.NewTimingWheel(time.Millisecond * 10)
	s := &Server{
		msgHandler: common.NewMsgHandle(100, 1024),
		connMgr:    common.NewConnManager(),
		tw:         tw,
	}
	return s
}

func (s *Server) Start() {
	go func() {
		s.msgHandler.StartWorkerPool()
		s.tw.Start()

		addr, err := net.ResolveTCPAddr("tcp4", ":3000")
		if err != nil {
			logger.Error("resolve tcp addr err: ", err)
			return
		}

		listener, err := net.ListenTCP("tcp4", addr)
		if err != nil {
			logger.Error("listen", "tcp4", "err", err)
			return
		}

		logger.Info("start server  aa", " success, now listenning...")

		cid := &atomic.Uint32{}
		cid.Set(0)

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				logger.Error("Accept err ", err)
				continue
			}

			logger.Info("Get conn remote addr = ", conn.RemoteAddr().String())
			if s.connMgr.Len() >= 12000 {
				conn.Close()
				continue
			}

			dealConn := common.NewConnection(common.Tcp, conn, cid.Get(), 1024, 4096)
			dealConn.SetTcpConn(conn)
			dealConn.SetOnClose(s.OnConnStop)

			var tm int64
			tm = 10
			ttl := time.Duration(tm)*time.Second
			s.tw.NewWheel(map[string]interface{}{"conn": dealConn, "ttl": tm}, ttl, s.beatHeart)

			s.connMgr.Add(dealConn)
			cid.Incr()

			go s.startReader(dealConn)
			go s.startWriter(dealConn)

			if s.OnConnStart != nil {
				s.OnConnStart(dealConn)
			}
		}
	}()
}

func (s *Server) Stop() {
	logger.Info("[STOP] server , name aa")

	s.tw.Stop()
	s.connMgr.ClearConn()
}

func (s *Server) AddRouter(msgId byte, router common.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

func (s *Server) GetConnMgr() *common.ConnManager {
	return s.connMgr
}

func (s *Server) SetOnConnStart(hookFunc func(*common.Connection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(*common.Connection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) beatHeart(param map[string]interface{}) error {
	conn := param["conn"].(*common.Connection)
	ttl := param["ttl"].(int64)
	t := time.Now().Unix() - conn.GetUpdateTime()

	if t > ttl {
		conn.Close()
	}
	return nil
}

func (s *Server) removeConnect(c *common.Connection) {
	logger.Info(c.RemoteAddr().String(), "[conn Reader exit!]")
	c.Close()
	s.connMgr.Remove(c)
}

func (s *Server) startReader(c *common.Connection) {
	logger.Info("[Reader Goroutine is running]")
	defer s.removeConnect(c)

	for {
		dp := newDataPack(4096)
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, headData); err != nil {
			logger.Error("server read msg head error ", err)
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
			if _, err := io.ReadFull(c.Conn, data); err != nil {
				logger.Error("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)

		req := common.Request{
			Conn: c,
			Msg:  msg,
		}

		go s.msgHandler.DoReceive(&req)
	}
}

func (s *Server) startWriter(c *common.Connection) {
	logger.Info("[Writer Goroutine is running]")
	defer logger.Info(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case msg, ok := <-c.GetMsgChan():
			if ok {
				dp := newDataPack(4096)
				data, err := dp.Pack(msg)
				if err != nil {
					logger.Error("Pack error msg id = ", msg.Id)
					return
				}
				if _, err := c.Conn.Write(data); err != nil {
					logger.Error("Send Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				logger.Info("msgBuffChan is Closed")
				break
			}
			//fmt.Printf("Send data succ! data = %+v\n", data)
		case msg, ok := <-c.GetMsgBuffChan():
			if ok {
				dp := newDataPack(4096)
				data, err := dp.Pack(msg)
				if err != nil {
					logger.Error("Pack error msg id = ", msg.Id)
					return
				}
				if _, err := c.Conn.Write(data); err != nil {
					logger.Error("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				logger.Info("msgBuffChan is Closed")
				break
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}