package scheduler

import (
	"encoding/json"
	"github.com/robfig/cron/v3"
	"sapi/pkg/client/etcdv3"
	"sapi/pkg/logger"
)

type Server struct {
	cron *cron.Cron
	client *etcdv3.Client

	registry  *Registry
	handles *Handlers
	jobs *Jobs
	works *Works
}

func NewServer(client *etcdv3.Client) *Server {
	return &Server{
		cron: cron.New(),
		client: client,
		handles: newHandlers(client),
		registry: newRegistry(client),
		jobs: newJobs(client),
		works: newWorks(client),
	}
}

func (s *Server) GetJobs() *Jobs {
	return s.jobs
}

func (s *Server) GetRegistry() *Registry {
	return s.registry
}

func (s *Server) GetWorks() *Works {
	return s.works
}

func (s *Server) Bootstrap() error {
	err := s.client.ConnectStatus()
	if err != nil {
		return err
	}

	jobs,err := s.jobs.ListJob()
	if err != nil {
		return err
	}

	for _, job := range jobs {
		s.addJobToCron(job)
	}
	s.cron.Start()

	wc, err := s.jobs.Watch()
	if err != nil {
		return err
	}

	go func() {
		for {
			ev, err := wc.Next()
			if err != nil {
				logger.Error(err)
				return
			}

			if ev.Type == etcdv3.KeyCreate {
				var job Job
				err = json.Unmarshal(ev.Value, &job)
				if err != nil {
					logger.Error(err)
					return
				}

				s.addJobToCron(&job)
			}
		}
	}()

	rwc, err := s.registry.Watch()
	if err != nil {
		return err
	}

	go func() {
		for {
			ev, err := rwc.Next()
			if err != nil {
				logger.Error(err)
				return
			}

			if ev.Type == etcdv3.KeyDelete {
				var clientInfo ClientInfo
				err = json.Unmarshal(ev.Value, &clientInfo)
				if err != nil {
					logger.Error(err)
					return
				}

				s.handles.DelClient(clientInfo.ID)
			}
		}
	}()
	return nil
}

func (s *Server) Stop() {
	s.client.Close()
	s.cron.Stop()
}

func (s *Server) addJobToCron(job *Job) {
	job.SetHandlers(s.handles)
	job.SetWorks(s.works)
	s.cron.AddJob(job.Cron, job)
}