package common

import (
	"github.com/26597925/EastCloud/pkg/logger"
	"strconv"
)

type MsgHandle struct {
	Apis           map[byte]IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan *Request
	MaxWorkerTaskLen uint32
}

func NewMsgHandle(workerPoolSize uint32, maxWorkerTaskLen uint32) *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[byte]IRouter),
		WorkerPoolSize: workerPoolSize,
		TaskQueue: make([]chan *Request, workerPoolSize),
		MaxWorkerTaskLen: maxWorkerTaskLen,
	}
}

func (mh *MsgHandle) DoReceive(request *Request) {
	if mh.WorkerPoolSize> 0 {
		mh.SendMsgToTaskQueue(request)
	} else {
		go mh.DoMsgHandler(request)
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request *Request) {
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	mh.TaskQueue[workerID] <- request
}

func (mh *MsgHandle) DoMsgHandler(request *Request) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		logger.Info("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	request.Conn.UpdateTime() // 当每次收到信息，更新链接时间
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId byte, router IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}

	mh.Apis[msgId] = router
}

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan *Request) {
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan *Request, mh.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}
