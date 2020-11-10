package router

import (
	"github.com/26597925/EastCloud/cmd/hello/controller"
	"github.com/26597925/EastCloud/pkg/server/api"
	"github.com/26597925/EastCloud/pkg/server/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/reflection"
)

func GrpcRouter(svr api.Server)  error {
	gr := svr.(*grpc.Server).Server
	helloworld.RegisterGreeterServer(gr, &controller.Greeter{})
	reflection.Register(gr)
	return nil
}
