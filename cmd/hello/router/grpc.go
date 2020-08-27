package router

import (
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
	"sapi/cmd/hello/controller"
	"sapi/pkg/server/api"
	"sapi/pkg/server/grpc"
)

func GrpcRouter(svr api.Server)  error {
	gr := svr.(*grpc.Server).Server
	helloworld.RegisterGreeterServer(gr, &controller.Greeter{})
	reflection.Register(gr)
	return nil
}
