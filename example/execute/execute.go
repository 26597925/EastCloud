package main

import (
	"github.com/google/uuid"
	"sapi/example/execute/handlers"
	"sapi/pkg/client/etcdv3"
	"sapi/pkg/logger"
	"sapi/pkg/scheduler"
)

func main() {
	cli := etcdv3.NewOptions().Build()
	sc := scheduler.NewClient(cli)
	sc.AddHandler(&handlers.HelloHandler{})
	sc.AddHandler(&handlers.DemoHandler{})
	sc.AddHandler(&handlers.TestHandler{})
	err := sc.Bootstrap()

	ss := scheduler.NewServer(cli)
	err = ss.Bootstrap()
	if err != nil {
		logger.Error(err)
		return
	}
	ss.GetJobs().AddJob(&scheduler.Job{
		ID:    uuid.New().String(),
		HandlerName: "HelloHandler",
		Cron:  "@every 5s",
		Param: "aa=bb&&cc=gg",
		Status: 1,
	})

	for {
		select {

		}
	}
}
