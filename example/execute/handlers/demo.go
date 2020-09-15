package handlers

import (
	"fmt"
	"sapi/pkg/scheduler"
)

type DemoHandler struct {

}

func (h *DemoHandler) GetNme() string {
	return "DemoHandler"
}

func (h *DemoHandler) Run(job *scheduler.Job) error {
	fmt.Println("参数:", job.Param)
	return nil
}
