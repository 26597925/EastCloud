package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/26597925/EastCloud/pkg/client/etcdv3"
)

const (
	NotExecuted = iota
	Executed
	Fail
	Finish
)

const (
	WorkPrefix = "/sapi/scheduler/work"
)

type Work struct {
	ID string
	JobId string
	ClientId string
	HandlerName string
	Time int64
	Status int
}

type Works struct {
	context context.Context
	client *etcdv3.Client
}

func newWorks(client *etcdv3.Client) *Works {
	return &Works{
		context: context.Background(),
		client: client,
	}
}

func (w *Works) putWork(work *Work) error {
	key := fmt.Sprintf("%s/%s", WorkPrefix, work.ID)

	data, err := json.Marshal(work)
	if err != nil {
		return err
	}

	_, err = w.client.Put(w.context, key, string(data))
	if err != nil {
		return err
	}

	return nil
}

func (w *Works) ClientWork(clientId string) ([]*Work, error){
	list, err := w.client.GetPrefix(w.context, WorkPrefix)
	if err != nil {
		return nil, err
	}

	var works []*Work
	for _, item := range list {
		var work Work
		err = json.Unmarshal(item, &work)
		if err != nil {
			return nil, err
		}

		works = append(works, &work)
	}

	return works, nil
}

func (w *Works) Watch() (*etcdv3.Watcher, error) {
	return w.client.NewWatchWithPrefixKey(WorkPrefix)
}