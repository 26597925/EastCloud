package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	retcd "github.com/coreos/etcd/clientv3"
	"io/ioutil"
	"sapi/pkg/registry"
	"sapi/pkg/server/api"
	"sync"
	"time"
)

type etcdRegistry struct {
	options *registry.Options

	client *retcd.Client
	lease  retcd.LeaseID
	register  sync.Map

	timeout time.Duration
}

func configure(e *etcdRegistry, opts *registry.Options) error {
	if opts.Prefix == "" {
		opts.Prefix = registry.DefaultPrefix
	}

	endpoints := opts.Endpoints
	if len(opts.Endpoints) == 0 {
		endpoints = []string{"localhost:2379"}
	}

	e.timeout = time.Duration(opts.Timeout) * time.Second
	if opts.Timeout == 0 {
		e.timeout = 3 * time.Second
	}

	config := retcd.Config{
		Endpoints:   endpoints,
		DialTimeout: e.timeout,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
	}

	if opts.BasicAuth {
		config.Username = opts.Username
		config.Password = opts.Password
	}

	tlsEnabled := false
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	if opts.CaCert != "" {
		certBytes, err := ioutil.ReadFile(opts.CaCert)
		if err != nil {
			return err
		}

		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(certBytes)

		if ok {
			tlsConfig.RootCAs = caCertPool
		}
		tlsEnabled = true
	}

	if opts.CertFile != "" && opts.KeyFile != "" {
		tlsCert, err := tls.LoadX509KeyPair(opts.CertFile, opts.KeyFile)
		if err != nil {
			return err
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
		tlsEnabled = true
	}

	if tlsEnabled {
		config.TLS = tlsConfig
	}

	client, err := retcd.New(config)
	if err != nil {
		return err
	}

	e.client = client
	return nil
}

func NewRegistry(opts *registry.Options) (registry.Registry, error) {
	e := &etcdRegistry{
		options: opts,
	}

	err := configure(e, opts)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (e *etcdRegistry) Register(opt api.Option) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(),  e.timeout)
	defer cancel()

	if e.lease > 0 {
		if _, err := e.client.KeepAliveOnce(context.TODO(), e.lease); err != nil {
			return err
		}
	}

	key := fmt.Sprintf("%s/%s/%s", e.options.Prefix, e.Type(), opt.GetName())
	service := &registry.Service{
		Driver:    opt.GetDriver(),
		Name:      opt.GetName(),
		ID:        opt.GetId(),
		Version:   opt.GetVersion(),
		Region:    opt.GetRegion(),
		Zone:      opt.GetZone(),
		GroupName: opt.GetGroupName(),
		IP:        opt.GetIP(),
		Port:      opt.GetPort(),
	}
	val, err := json.Marshal(service)

	var lgr *retcd.LeaseGrantResponse
	if e.options.TTL > 0 {
		ttl := time.Duration(e.options.TTL) * time.Second
		lgr, err = e.client.Grant(ctx, int64(ttl.Seconds()))
		if err != nil {
			return err
		}
		e.lease = lgr.ID
	}

	var opOptions []retcd.OpOption
	if e.lease != 0 {
		opOptions = append(opOptions, retcd.WithLease(e.lease))
	}

	if _, err = e.client.Put(ctx, key, string(val), opOptions...); err != nil {
		return err
	}

	e.register.Store(key, val)
	return nil
}

func (e *etcdRegistry) Deregister(sv *registry.Service) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	key := fmt.Sprintf("%s/%s/%s", e.options.Prefix, e.Type(), sv.Name)
	_, err := e.client.Delete(ctx, key)
	if err == nil {
		e.register.Delete(key)
	}
	return err
}

func (e *etcdRegistry) GetService(name string) (*registry.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	key := fmt.Sprintf("%s/%s/%s", e.options.Prefix, e.Type(), name)

	rsp, err := e.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if rsp == nil || len(rsp.Kvs) == 0 {
		return nil, fmt.Errorf("source not found: %s", key)
	}

	var s *registry.Service
	err = json.Unmarshal(rsp.Kvs[0].Value, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (e *etcdRegistry) ListServices() ([]*registry.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	rsp, err := e.client.Get(ctx, e.options.Prefix, retcd.WithPrefix(), retcd.WithSerializable())
	if err != nil {
		return nil, err
	}
	if len(rsp.Kvs) == 0 {
		return []*registry.Service{}, nil
	}

	services := make([]*registry.Service, 0, len(rsp.Kvs))
	for _, n := range rsp.Kvs {
		var s *registry.Service
		err = json.Unmarshal(n.Value, &s)

		if err != nil {
			continue
		}

		services = append(services, s)
	}

	return services, nil
}

func (e *etcdRegistry) Watch() (registry.Watcher, error) {
	return newEtcdWatcher(e)
}

func (e *etcdRegistry) Close() (err error) {
	var wg sync.WaitGroup
	e.register.Range(func(k, v interface{}) bool {
		wg.Add(1)
		go func(v interface{}) {
			defer wg.Done()
			var s *registry.Service
			err = json.Unmarshal(v.([]byte), &s)
			if err == nil {
				err = e.Deregister(s)
			}
		}(v)
		return true
	})
	wg.Wait()

	if e.lease > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err = e.client.Revoke(ctx, e.lease)
		cancel()
		return err
	}
	return nil
}

func (e *etcdRegistry) Type() string {
	return "etcdv3"
}