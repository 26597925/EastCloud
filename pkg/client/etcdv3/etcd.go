package etcdv3

import (
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
)

type Client struct {
	opt 		*Options
	config      clientv3.Config
	*clientv3.Client
}

func newClient(options *Options) (*Client, error) {
	config := options.BuildConfig()
	client, err := clientv3.New(config)

	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), config.DialTimeout)
	defer cancel()
	_, err = client.Status(timeoutCtx, config.Endpoints[0])
	if err != nil {
		return nil, err
	}

	return &Client{
		opt:		 options,
		config:		 config,
		Client:      client,
	}, nil
}

func (client *Client) GetValue(ctx context.Context, key string) ([]byte,  error) {
	rp, err := client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(rp.Kvs) > 0 {
		return rp.Kvs[0].Value, nil
	}

	return	nil, errors.New("no data")
}

func (client *Client) GetPrefix(ctx context.Context, prefix string) (map[string][]byte, error) {
	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	vars := make(map[string][]byte)
	for _, kv := range resp.Kvs {
		vars[string(kv.Key)] = kv.Value
	}

	return vars, nil
}

func (client *Client) GetPrefixLimit(ctx context.Context, prefix string, limit int64) (map[string][]byte, error) {
	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithLimit(limit))
	if err != nil {
		return nil, err
	}

	vars := make(map[string][]byte)
	for _, kv := range resp.Kvs {
		vars[string(kv.Key)] = kv.Value
	}

	return vars, nil
}

// put a key not exist
func (client *Client) PutNotExist(ctx context.Context, key, value string) ([]byte, error) {
	txn := client.Txn(ctx)

	txnResp, err := txn.If(clientv3.Compare(clientv3.Version(key), "=", 0)).
		Then(clientv3.OpPut(key, value)).
		Else(clientv3.OpGet(key)).
		Commit()

	if err != nil {
		return nil, err
	}

	if !txnResp.Succeeded {
		return nil, errors.New("put val is fail")
	}

	return txnResp.Responses[0].GetResponseRange().Kvs[0].Value, nil
}

func (client *Client) Update(ctx context.Context, key, value, oldValue string) error {
	txn := client.Txn(ctx)

	txnResp, err := txn.If(clientv3.Compare(clientv3.Value(key), "=", oldValue)).
		Then(clientv3.OpPut(key, value)).
		Commit()

	if err != nil {
		return err
	}

	if !txnResp.Succeeded {
		return errors.New("update val is fail")
	}

	return nil
}

func (client *Client) NewWatch(key string) (*Watcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	wc := clientv3.NewWatcher(client.Client)
	exit := make(chan bool, 1)

	w := &Watcher{
		ch:          make(chan *Event),
		exit:        exit,
	}

	go func() {
		<-exit
		cancel()
	}()

	ch := wc.Watch(ctx, key)

	go w.run(wc, ch)

	return w, nil
}

func (client *Client) NewWatchWithPrefixKey(prefixKey string) (*Watcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	wc := clientv3.NewWatcher(client.Client)
	exit := make(chan bool, 1)

	w := &Watcher{
		ch:          make(chan *Event),
		exit:        exit,
	}

	go func() {
		<-exit
		cancel()
	}()

	ch := wc.Watch(ctx, prefixKey, clientv3.WithPrefix())

	go w.run(wc, ch)

	return w, nil
}

func (client *Client) NewRegistry() *Registry {
	r := &Registry{
		client: client,
		timeout: client.config.DialTimeout,
		register: make(map[string]int),
		leases:   make(map[string]clientv3.LeaseID),
	}

	return r
}

func (client *Client) Transfer(from string, to string, value string) (bool, error) {
	txnResponse, err := client.Txn(context.Background()).If(
		clientv3.Compare(clientv3.Value(from), "=", value)).
		Then(
			clientv3.OpDelete(from),
			clientv3.OpPut(to, value),
		).Commit()

	if err != nil {
		return false, err
	}

	return txnResponse.Succeeded, nil
}