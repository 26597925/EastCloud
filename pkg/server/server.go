package server

import (
	"github.com/26597925/EastCloud/pkg/server/api"
	"github.com/26597925/EastCloud/pkg/server/gin"
	"github.com/26597925/EastCloud/pkg/server/grpc"
)

type Option func(*Options)

func NewServer(option api.Option) api.Server {
	if option.GetDriver() == "gin" {
		return gin.NewServer(option)
	} else if option.GetDriver() == "grpc" {
		return grpc.NewServer(option)
	}

	return nil
}