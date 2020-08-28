package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"log"
	"os"
	"sapi/ph/grpc_server/test"
	gCliect "sapi/pkg/client/grpc"
)

func test_client() {
	//建立连接到gRPC服务
	conn, err := grpc.Dial("127.0.0.1:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 函数结束时关闭连接
	defer conn.Close()

	// 创建Waiter服务的客户端
	t := test.NewWaiterClient(conn)

	// 模拟请求数据
	req := "test"
	// os.Args[1] 为用户执行输入的参数 如：go run ***.go 123
	if len(os.Args) > 1 {
		req = os.Args[1]
	}

	// 调用gRPC接口
	tr, err := t.DoMD5(context.Background(), &test.Req{JsonStr: req})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("服务端响应: %s", tr.BackJson)
}

func main() {
	//for i := 0; i <= 1000; i++ {
	//	test_client()
	//}

	res := &test.Res{}
	cc := gCliect.NewOptions(gCliect.Address("127.0.0.1:8000"), gCliect.GrpcDialOptions(grpc.WithInsecure())).Build()

	for i := 0; i <= 1; i++ {
		req := fmt.Sprintf("test%d", i)
		err := cc.Call(context.Background(), "/test.Waiter/DoMD5", &test.Req{JsonStr: req}, res)
		fmt.Println(err)
		data := fmt.Sprintf("number:%d    res:%v", i, res)
		fmt.Println(data)
	}

	cc1 := gCliect.NewOptions(gCliect.Address("127.0.0.1:8001"), gCliect.GrpcDialOptions(grpc.WithInsecure())).Build()
	res1 := &helloworld.HelloReply{}
	err := cc1.Call(context.Background(), "/helloworld.Greeter/SayHello", &helloworld.HelloRequest{
		Name: "World",
	}, res1)

	fmt.Println(err)
	fmt.Println(res1.Message)
}