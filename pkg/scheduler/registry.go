package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/26597925/EastCloud/pkg/client/etcdv3"
)

const (
	RegistryPrefix = "/sapi/scheduler/client"
)

type Handler interface {
	GetNme() string
	Run(job *Job) error
}

type Registry struct {
	ttl int64
	client    *etcdv3.Client
	registry  *etcdv3.Registry
}

func newRegistry(client *etcdv3.Client) *Registry {
	return &Registry{
		ttl: 6000,
		client: client,
		registry: client.NewRegistry(),
	}
}

func (r *Registry) Register(client *ClientInfo) error {
	key := fmt.Sprintf("%s/%s", RegistryPrefix, client.ID)

	data, err := json.Marshal(client)
	if err != nil {
		return err
	}

	err = r.registry.Register(key, string(data), r.ttl)
	if err != nil {
		return err
	}

	return nil
}

func (r *Registry) Deregister(id string)  error {
	key := fmt.Sprintf("%s/%s", RegistryPrefix, id)
	err := r.registry.Deregister(key)

	return err
}

func (r *Registry) GetClient(id string) (*ClientInfo, error) {
	key := fmt.Sprintf("%s/%s", RegistryPrefix, id)
	data, err := r.client.GetValue(context.Background(), key)
	if err != nil {
		return nil, err
	}

	var client ClientInfo
	err = json.Unmarshal(data, &client)
	if err != nil {
		return nil, err
	}

	return &client, nil
}


func (r *Registry) ListClients() ([]*ClientInfo, error) {
	list, err := r.client.GetPrefix(context.Background(), RegistryPrefix)
	if err != nil {
		return nil, err
	}

	clients := make([]*ClientInfo, 0, len(list))
	for _, data :=range list {
		var client ClientInfo
		err = json.Unmarshal(data, &client)

		if err != nil {
			continue
		}

		clients = append(clients, &client)
	}

	return clients, nil
}

func (r *Registry) Watch() (*etcdv3.Watcher, error) {
	return r.client.NewWatchWithPrefixKey(RegistryPrefix)
}