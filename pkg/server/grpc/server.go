package grpc

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"sapi/pkg/server/api"
)

type Server struct {
	*grpc.Server
	listener net.Listener
	option api.Option
}
//https://github.com/fullstorydev/grpcui
func NewServer(option api.Option) api.Server {
	return &Server{
		Server: grpc.NewServer(option.GetGrpcOptions()...),
		option: option,
	}
}

func (svr *Server) Handler(handler api.Handler) error {
	if handler == nil {
		return errors.New("handler errors is nil ")
	}

	return handler(svr)
}

func (svr *Server) Init() error {
	addr := fmt.Sprintf("%s:%d",  svr.option.GetIP(), svr.option.GetPort())
	listener, err := net.Listen("tcp", addr)
	svr.listener = listener
	return err
}

func (svr *Server) Start() error {
	return svr.Server.Serve(svr.listener)
}

func (svr *Server) Stop() error {
	svr.Server.Stop()
	return nil
}

func (svr *Server) GracefulStop(ctx context.Context) error {
	svr.Server.GracefulStop()
	return nil
}

func (svr *Server) GetOption() api.Option {
	return svr.option
}