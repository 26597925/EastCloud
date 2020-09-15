package handlers

import (
	"fmt"
	"sapi/pkg/scheduler"
)

type TestHandler struct {

}

func (h *TestHandler) GetNme() string {
	return "TestHandler"
}

func (h *TestHandler) Run(job *scheduler.Job) error {
	fmt.Println("参数:", job.Param)
	return nil
}
