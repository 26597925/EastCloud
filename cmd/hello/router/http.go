package router

import (
	"github.com/26597925/EastCloud/cmd/hello/controller"
	"github.com/26597925/EastCloud/pkg/server/api"
	"github.com/26597925/EastCloud/pkg/server/gin"
)

func HttpRouter(svr api.Server)  error {
	router := svr.(*gin.Server)
	demo := &controller.Demo{}
	router.GET("/welcome", demo.Query)
	router.GET("/list", demo.List)
	return nil
}