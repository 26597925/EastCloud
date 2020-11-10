package agent

import (
	"github.com/26597925/EastCloud/internal/agent/msg"
	"github.com/26597925/EastCloud/internal/agent/router/tcpserver"
	"github.com/26597925/EastCloud/internal/agent/router/web"
	"github.com/26597925/EastCloud/pkg/network/common"
	"github.com/26597925/EastCloud/pkg/network/tcp"
	"github.com/26597925/EastCloud/pkg/network/websocket"
)

type Server struct {
	web *websocket.Server
	server *tcp.Server
}

func NewServer() *Server {
	svr := &Server{
		web:     websocket.NewWebsocket(),
		server:  tcp.NewServer(),
	}

	return svr
}

func (s *Server) GetTcpConnectMng() *common.ConnManager {
	return s.server.GetConnMgr()
}

func (s *Server) Start() {
	s.server.AddRouter(msg.Heartbeat, &tcpserver.HeartbeatRouter{})
	s.server.AddRouter(msg.OK, &tcpserver.OkRouter{
		Web: s.web,
	})
	s.server.AddRouter(msg.ERROR, &tcpserver.ErrorRouter{
		Web: s.web,
	})
	s.server.Start()

	s.web.AddRouter(msg.Command, &web.CommandRouter{
		Server: s.server,
	})
	s.web.AddRouter(msg.File, &web.FileRouter{
		Server: s.server,
	})
	s.web.Start()
}

func (s *Server) Stop() {
	s.server.Stop()
	s.web.Stop()
}