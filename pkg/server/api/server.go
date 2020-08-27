package api

import (
	"context"
	"google.golang.org/grpc"
)

type Server interface {
	Init() error
	Handler(Handler) error
	Start() error
	Stop() error
	GracefulStop(ctx context.Context) error
	GetOption() Option
}

type Option interface {
	GetDriver() string
	GetName() string
	GetVersion() string
	GetId() string
	GetRegion() string
	GetZone() string
	GetGroupName() string
	GetIP() string
	GetPort() int
	GetGrpcOptions() []grpc.ServerOption
}

type Handler func(Server) error