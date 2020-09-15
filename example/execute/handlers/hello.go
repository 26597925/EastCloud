package handlers

import (
	"sapi/pkg/logger"
	"sapi/pkg/scheduler"
)

type HelloHandler struct {

}

func (h *HelloHandler) GetNme() string {
	return "HelloHandler"
}

func (h *HelloHandler) Run(job *scheduler.Job) error {
	logger.InfoF("%v", job)
	return nil
}
