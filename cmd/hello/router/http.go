package router

import (
	"sapi/cmd/hello/controller"
	"sapi/pkg/server/api"
	"sapi/pkg/server/gin"
)

func HttpRouter(svr api.Server)  error {
	router := svr.(*gin.Server)
	demo := &controller.Demo{}
	router.GET("/welcome", demo.Query)
	router.GET("/list", demo.List)
	return nil
}