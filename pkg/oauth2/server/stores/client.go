package stores

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/26597925/EastCloud/pkg/store"
	"time"
)

var ClientPrefix = "/sapi/oauth2/client"

type ClientStore struct {
	Store store.Store
	Prefix string
	Expire time.Duration
}

type Client struct {
	ID     string
	Secret string
	Domain string
	UserID string
}

func (cs *ClientStore) GetClient(ctx context.Context, id string) (*Client, error) {
	prefix := cs.Prefix
	if prefix == "" {
		prefix = ClientPrefix
	}
	key := fmt.Sprintf("%s/%s", prefix, id)

	d, err := cs.Store.Get(ctx, key)
	if err != nil {
		return nil, errors.New("not found")
	}

	var c Client
	err = json.Unmarshal([]byte(d), &c)
	if err != nil {
		return nil, errors.New("data parse errors")
	}

	return &c, nil
}

func (cs *ClientStore) Set(ctx context.Context, id string, cli *Client)  error {
	prefix := cs.Prefix
	if prefix == "" {
		prefix = ClientPrefix
	}
	key := fmt.Sprintf("%s/%s", prefix, id)

	b, err := json.Marshal(cli)
	if err != nil {
		return err
	}

	err = cs.Store.Set(ctx, key, string(b), cs.Expire)
	if err != nil {
		return err
	}

	return nil
}