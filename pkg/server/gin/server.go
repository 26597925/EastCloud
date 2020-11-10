package gin

import (
	"context"
	"errors"
	"fmt"
	"github.com/26597925/EastCloud/pkg/server/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	*gin.Engine
	*http.Server
	option api.Option
}

func NewServer(option api.Option) api.Server {
	return &Server{
		Engine: gin.New(),
		option: option,
	}
}

func (svr *Server) Handler(handler api.Handler) error {
	if handler == nil {
		return errors.New("handler errors is nil ")
	}

	return handler(svr)
}

func (svr *Server) Init() error {
	addr := fmt.Sprintf("%s:%d",  svr.option.GetIP(), svr.option.GetPort())
	svr.Server = &http.Server{
		Addr:           addr,
		Handler:        svr.Engine,
	}
	return nil
}

func (svr *Server) Start() error {
	return svr.Server.ListenAndServe()
}

func (svr *Server) Stop() error {
	return svr.Server.Close()
}

func (svr *Server) GracefulStop(ctx context.Context) error {
	return svr.Server.Shutdown(ctx)
}

func (svr *Server) GetOption() api.Option {
	return svr.option
}