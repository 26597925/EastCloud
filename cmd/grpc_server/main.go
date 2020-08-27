package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"google.golang.org/grpc/reflection"
	"sapi/ph/grpc_server/test"
	"sapi/pkg/server"
	"sapi/pkg/server/api"
	"sapi/pkg/server/grpc"
)

type testServer struct{}

//https://github.com/grpc-ecosystem/go-grpc-middleware
// 为server定义 DoMD5 方法 内部处理请求并返回结果
// 参数 (context.Context[固定], *test.Req[相应接口定义的请求参数])
// 返回 (*test.Res[相应接口定义的返回参数，必须用指针], error)
func (s *testServer) DoMD5(ctx context.Context, in *test.Req) (*test.Res, error) {
	fmt.Println("MD5方法请求JSON:"+in.JsonStr)
	return &test.Res{BackJson: "MD5 :" + fmt.Sprintf("%x", md5.Sum([]byte(in.JsonStr)))}, nil
}

func RegisterService(svr api.Server)  error{
	gr := svr.(*grpc.Server).Server
	test.RegisterWaiterServer(gr, &testServer{})
	reflection.Register(gr)
	return nil
}

func main() {

	srv := server.NewServer(server.NewOptions(server.Driver("grpc")))
	srv.Init()
	srv.Handler(RegisterService)
	srv.Start()

}
