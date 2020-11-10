package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/26597925/EastCloud/pkg/client/etcdv3"
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/google/uuid"
	"time"
)

const (
	JobPrefix = "/sapi/scheduler/job"
)

type Job struct {
	ID string
	Name string
	Remark string
	Cron string
	RunMode int
	ExecMode int
	HandlerName string
	Scheduler int
	ChildJobId string
	Timeout int64
	RetryNumber int
	UserName string
	Email string
	Param string
	Status int //是否开启

	works *Works
	handlers *Handlers
}

func (j *Job) SetWorks(works *Works) {
	j.works = works
}

func (j *Job) SetHandlers(handlers *Handlers) {
	j.handlers = handlers
}

func (j *Job) Run() {
	if j.Status == 1 {
		cls, err := j.handlers.FindClients(j.HandlerName)
		if err != nil {//写入错误快照，状态更改为失败
			logger.Error(err)
			return
		}

		cliId, err := SelectClient(j.Scheduler , cls)
		if err != nil {
			logger.Error(err)
			return
		}

		err = j.works.putWork(&Work{
			ID: uuid.New().String(),
			JobId: j.ID,
			HandlerName: j.HandlerName,
			Time: time.Now().Unix(),
			Status: NotExecuted,
			ClientId: cliId,
		})
		if err != nil {//写入错误快照，状态更改为失败
			logger.Error(err)
			return
		}
	}
}

type Jobs struct {
	context context.Context
	client *etcdv3.Client
}

func newJobs(client *etcdv3.Client) *Jobs {
	return &Jobs{
		context: context.Background(),
		client: client,
	}
}

func (j *Jobs) GetFirstJobKey() (string, error) {
	return j.client.GetFirstKey(j.context, JobPrefix)
}

func (j *Jobs) CountJob() (int64, error) {
	return j.client.GetCount(j.context, JobPrefix)
}

func (j *Jobs) AddJob(job *Job) error {
	key := fmt.Sprintf("%s/%s", JobPrefix, job.ID)

	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	_, err = j.client.Put(j.context, key, string(data))
	if err != nil {
		return err
	}

	return nil
}

func (j *Jobs) GetJob(id string) (*Job, error) {
	key := fmt.Sprintf("%s/%s", JobPrefix, id)
	data, err := j.client.GetValue(j.context, key)
	if err != nil {
		return nil, err
	}

	var job Job
	err = json.Unmarshal(data, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (j *Jobs) ListJob() ([]*Job, error) {
	list, err := j.client.GetPrefix(j.context, JobPrefix)
	if err != nil {
		return nil, err
	}

	jobs := make([]*Job, 0, len(list))
	for _, data :=range list {
		var job Job
		err = json.Unmarshal(data, &job)

		if err != nil {
			continue
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (j *Jobs) PageJob(firstKey string, limit int64) ([]*Job, error) {
	list, err := j.client.GetPrefixLimit(j.context, firstKey, limit)
	if err != nil {
		return nil, err
	}

	jobs := make([]*Job, 0, len(list))
	for _, data :=range list {
		var job Job
		err = json.Unmarshal(data, &job)

		if err != nil {
			continue
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (j *Jobs) RemoveJob(id string) error {
	key := fmt.Sprintf("%s/%s", JobPrefix, id)
	_, err := j.client.Delete(j.context, key)
	if err != nil {
		return err
	}

	return nil
}

func (j *Jobs) Watch() (*etcdv3.Watcher, error) {
	return j.client.NewWatchWithPrefixKey(JobPrefix)
}