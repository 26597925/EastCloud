package etcdv3

import (
	"crypto/tls"
	"github.com/26597925/EastCloud/pkg/util/crypto"
	"github.com/coreos/etcd/clientv3"
	"time"
)

type Options struct {
	Endpoints []string
	Timeout   int

	BasicAuth bool
	Username  string
	Password  string

	CertFile  string
	KeyFile   string
	CaCert    string
}

type Option func(*Options)

func NewOptions(opts ...Option) *Options {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}

	return options
}

func (opt *Options) BuildConfig () clientv3.Config {
	endpoints := opt.Endpoints
	if len(opt.Endpoints) == 0 {
		endpoints = []string{"localhost:2379"}
	}

	dialTimeout := time.Duration(opt.Timeout) * time.Second
	if opt.Timeout == 0 {
		dialTimeout = 3 * time.Second
	}

	config := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
	}

	if opt.BasicAuth {
		config.Username = opt.Username
		config.Password = opt.Password
	}

	tlsEnabled := false
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	caCertPool, err := crypto.ReadPEM(opt.CaCert)
	if err == nil {
		tlsConfig.RootCAs = caCertPool
		tlsEnabled = true
	}

	tsl, err := crypto.ReadTls(opt.CertFile, opt.KeyFile)
	if err == nil {
		tlsConfig.Certificates = tsl
		tlsEnabled = true
	}

	if tlsEnabled {
		config.TLS = tlsConfig
	}

	return config
}

func (opt *Options) Build() *Client {
	cc, err := newClient(opt)

	if err != nil {
		return nil
	}

	return cc
}