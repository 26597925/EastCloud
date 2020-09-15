package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sapi/pkg/client/etcdv3"
	"sort"
)

const (
	HandlerPrefix = "/sapi/scheduler/handler"
)

type Handlers struct {
	context context.Context
	client    *etcdv3.Client
}

type HandlerInfo struct {
	Name string
	Clients []string
}

func newHandlers(client *etcdv3.Client) *Handlers {
	return &Handlers{
		context: context.Background(),
		client: client,
	}
}

func (h *Handlers) PutHandler(clientId, name string) error {
	key := fmt.Sprintf("%s/%s", HandlerPrefix, name)

	rp, err := h.client.Get(h.context, key)
	if err != nil {
		return err
	}

	var handlerInfo HandlerInfo
	if len(rp.Kvs) > 0 {
		err = json.Unmarshal(rp.Kvs[0].Value, &handlerInfo)
		if err != nil {
			return err
		}
		clients := handlerInfo.Clients
		sort.Strings(clients)
		index := sort.SearchStrings(clients, clientId)
		if index == len(clients) {
			clients = append(clients, clientId)
		} else {
			if clients[index] != clientId {
				clients = append(clients, clientId)
			}
		}
		handlerInfo.Clients = clients
	} else {
		handlerInfo = HandlerInfo{
			Name: name,
			Clients: []string{clientId},
		}
	}

	data, err := json.Marshal(handlerInfo)
	if err != nil {
		return err
	}

	_, err = h.client.Put(h.context, key, string(data))
	if err != nil {
		return err
	}
	return nil
}

func (h *Handlers) ListHandlers() ([]string, error) {
	list, err := h.client.GetPrefix(h.context, HandlerPrefix)
	if err != nil {
		return nil, err
	}

	handlers := make([]string, 0, len(list))
	for _, data := range list {
		var handlerInfo HandlerInfo
		err = json.Unmarshal(data, &handlerInfo)

		if err != nil {
			continue
		}

		handlers = append(handlers, handlerInfo.Name)
	}

	return handlers, nil
}

func (h *Handlers) FindClients(name string) ([]string ,error) {
	key := fmt.Sprintf("%s/%s", HandlerPrefix, name)
	data, err := h.client.GetValue(h.context, key)
	if err != nil {
		return nil, err
	}

	var handlerInfo HandlerInfo
	err = json.Unmarshal(data, &handlerInfo)
	if err != nil {
		return nil, err
	}

	return handlerInfo.Clients, nil
}

//监听到client掉线后删除handler里面的客户端
func (h *Handlers) DelClient(clientId string) error {
	list, err := h.client.GetPrefix(h.context, HandlerPrefix)
	if err != nil {
		return err
	}

	for _, data := range list {
		var handlerInfo HandlerInfo
		err = json.Unmarshal(data, &handlerInfo)

		if err != nil {
			continue
		}

		clients := handlerInfo.Clients
		sort.Strings(clients)
		index := sort.SearchStrings(clients, clientId)
		if index < len(clients) {
			if clients[index] == clientId {
				clients = append(clients [:index], clients [index+1:]...)
				handlerInfo.Clients = clients

				data, err := json.Marshal(handlerInfo)
				if err != nil {
					return err
				}

				key := fmt.Sprintf("%s/%s", HandlerPrefix, handlerInfo.Name)
				_, err = h.client.Put(h.context, key, string(data))
				if err != nil {
					return err
				}
			}
		}

	}

	return nil
}