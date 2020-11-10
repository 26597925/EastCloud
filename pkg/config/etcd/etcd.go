package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/26597925/EastCloud/pkg/config/api"
	"io/ioutil"
	"time"

	cetcd "github.com/coreos/etcd/clientv3"
)

var (
	DefaultPrefix = "/sapi/config"
)
//https://github.com/mistaker/etcdTool
type Options struct {
	Name string

	Endpoints []string
	Timeout   int
	Format    string
	Prefix    string
	Debug     bool

	BasicAuth bool
	Username  string
	Password  string

	CertFile  string
	KeyFile   string
	CaCert    string
}

type Etcd struct {
	debug		bool
	prefix      string
	format      string
	config      cetcd.Config
	client      *cetcd.Client
}

func NewEtcd(options *Options) (*Etcd, error){
	endpoints := options.Endpoints
	if len(options.Endpoints) == 0 {
		endpoints = []string{"localhost:2379"}
	}

	dialTimeout := time.Duration(0)
	if options.Timeout != 0 {
		dialTimeout = time.Duration(options.Timeout) * time.Second
	}

	config := cetcd.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
	}

	if options.BasicAuth {
		config.Username = options.Username
		config.Password = options.Password
	}

	tlsEnabled := false
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	if options.CaCert != "" {
		certBytes, err := ioutil.ReadFile(options.CaCert)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(certBytes)

		if ok {
			tlsConfig.RootCAs = caCertPool
		}
		tlsEnabled = true
	}

	if options.CertFile != "" && options.KeyFile != "" {
		tlsCert, err := tls.LoadX509KeyPair(options.CertFile, options.KeyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
		tlsEnabled = true
	}

	if tlsEnabled {
		config.TLS = tlsConfig
	}

	client, err := cetcd.New(config)
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	_, err = client.Status(timeoutCtx, config.Endpoints[0])
	if err != nil {
		return nil, err
	}

	format := options.Format
	if format == "" {
		format = "yml"
	}

	mode := "pro"
	if options.Debug {
		mode = "dev"
	}

	prefix := options.Prefix
	if prefix == "" {
		prefix = DefaultPrefix
	}

	prefix = fmt.Sprintf("%v/%v/%v/%v", prefix, options.Name, mode, format)

	return &Etcd{
		debug:       options.Debug,
		prefix:      prefix,
		format: 	 format,
		config:		 config,
		client:      client,
	}, nil
}

func (c *Etcd) Watch() (api.Watcher, error) {
	return newWatcher(c, c.client.Watcher)
}

func (c *Etcd) Put(val string) {
	c.client.Put(context.Background(), c.prefix, val)
}

func (c *Etcd) Del(ctx context.Context, prefix string) (deleted int64, err error) {
	resp, err := c.client.Delete(ctx, prefix)
	if err != nil {
		return 0, err
	}
	return resp.Deleted, err
}

func (c *Etcd) Read(data interface{}) error {
	rsp, err := c.client.Get(context.Background(), c.prefix)
	if err != nil {
		return err
	}

	if rsp == nil || len(rsp.Kvs) == 0 {
		return fmt.Errorf("source not found: %s", c.prefix)
	}

	return api.Encoders[c.format].Decode(rsp.Kvs[0].Value, data)
}