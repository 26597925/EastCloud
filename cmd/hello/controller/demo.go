package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"sapi/cmd/hello/boot/engine"
	gCliect "sapi/pkg/client/grpc"
	"sapi/pkg/logger"
)

type Demo struct {
}

func (d *Demo) Query(c *gin.Context) {
	name, suc := c.GetQuery("name")
	if !suc {
		name = "World"
	}

	cc1 := gCliect.NewOptions(gCliect.Address("127.0.0.1:8089"), gCliect.GrpcDialOptions(grpc.WithInsecure())).Build()
	res1 := &helloworld.HelloReply{}
	err := cc1.Call(context.Background(), "/helloworld.Greeter/SayHello", &helloworld.HelloRequest{
		Name: name,
	}, res1)

	fmt.Println(err)
	c.Writer.WriteString(res1.Message)
}

func (d *Demo) List(c *gin.Context) {
	rgy := engine.GetServiceContext().GetRegistry()
	svrs, err := rgy.ListServices()
	logger.Info(err)
	logger.Info(svrs)

	for _, svr := range svrs {
		logger.Info(svr)
	}
}