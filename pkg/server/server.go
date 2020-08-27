package server

import (
	"sapi/pkg/server/api"
	"sapi/pkg/server/gin"
	"sapi/pkg/server/grpc"
)

type Option func(*Options)

func NewServer(option api.Option) api.Server{
	if option.GetDriver() == "gin" {
		return gin.NewServer(option)
	} else if option.GetDriver() == "grpc" {
		return grpc.NewServer(option)
	}

	return nil
}