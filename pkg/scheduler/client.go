package scheduler

import (
	"encoding/json"
	"github.com/26597925/EastCloud/pkg/client/etcdv3"
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/google/uuid"
)

type ClientInfo struct {
	ID      string
	Enable  bool
}

type Client struct {
	id  string
	ttl int64
	handlers map[string]Handler

	client *etcdv3.Client
	registry  *Registry
	handler *Handlers
	jobs *Jobs
	works *Works
}

func NewClient(client *etcdv3.Client) *Client {
	id := uuid.New().String()
	return &Client{
		id:       id,
		handlers: make(map[string]Handler),
		client: client,
		registry: newRegistry(client),
		handler: newHandlers(client),
		jobs: newJobs(client),
		works: newWorks(client),
	}
}

func (c *Client) AddHandler(handler Handler) {
	c.handlers[handler.GetNme()] = handler
}

func (c *Client) Bootstrap() error {
	err := c.client.ConnectStatus()
	if err != nil {
		return err
	}

	for _, handler := range c.handlers {
		c.handler.PutHandler(c.id, handler.GetNme())
	}

	client := &ClientInfo{
		ID:       c.id,
	}

	err = c.registry.Register(client)
	if err != nil {
		return err
	}

	wc, err := c.works.Watch()
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
				var work Work
				err = json.Unmarshal(ev.Value, &work)
				if err != nil {
					logger.Error(err)
					return
				}

				if work.ClientId == c.id && work.Status == NotExecuted {
					job, err := c.jobs.GetJob(work.JobId)
					if err != nil {
						logger.Error(err)
						return
					}

					go func() {
						work.Status = Executed
						c.works.putWork(&work)
						err := c.handlers[work.HandlerName].Run(job)
						if err == nil {
							work.Status = Finish
						} else {
							work.Status = Fail
						}
						c.works.putWork(&work)
					}()
				}
			}
		}
	}()
	return nil
}

func (c *Client) Stop() {
	c.client.Close()
}