package rabbitmq

import (
	"crypto/tls"
	"github.com/streadway/amqp"
	"regexp"
	"sapi/pkg/broker"
	"strings"
	"sync"
	"time"
)

type rabbitMQConn struct {
	opts 			*broker.Options

	Connection      *amqp.Connection
	Channel         *rabbitMQChannel
	ExchangeChannel *rabbitMQChannel

	sync.Mutex
	connected bool
	close     chan bool

	waitConnection chan struct{}
}

func newRabbitMQConn(opts *broker.Options) *rabbitMQConn {
	ret := &rabbitMQConn{
		opts: 			opts,
		close:          make(chan bool),
		waitConnection: make(chan struct{}),
	}
	// its bad case of nil == waitConnection, so close it at start
	close(ret.waitConnection)
	return ret
}

func (r *rabbitMQConn) Connect(secure bool, config *amqp.Config) error {
	r.Lock()

	// already connected
	if r.connected {
		r.Unlock()
		return nil
	}

	// check it was closed
	select {
	case <-r.close:
		r.close = make(chan bool)
	default:
		// no op
		// new conn
	}

	r.Unlock()

	return r.connect(secure, config)
}

func (r *rabbitMQConn) Publish(exchange, key string, msg amqp.Publishing) error {
	return r.ExchangeChannel.Publish(exchange, key, msg)
}

func (r *rabbitMQConn) Consume(queue, key string, headers amqp.Table, qArgs amqp.Table, autoAck, durableQueue bool) (*rabbitMQChannel, <-chan amqp.Delivery, error) {
	consumerChannel, err := newRabbitChannel(r.Connection, r.opts.RabbitMq.PrefetchCount, r.opts.RabbitMq.PrefetchGlobal)
	if err != nil {
		return nil, nil, err
	}

	if durableQueue {
		err = consumerChannel.DeclareDurableQueue(queue, qArgs)
	} else {
		err = consumerChannel.DeclareQueue(queue, qArgs)
	}

	if err != nil {
		return nil, nil, err
	}

	deliveries, err := consumerChannel.ConsumeQueue(queue, autoAck)
	if err != nil {
		return nil, nil, err
	}

	err = consumerChannel.BindQueue(queue, key, r.opts.RabbitMq.ExchangeKey, headers)
	if err != nil {
		return nil, nil, err
	}

	return consumerChannel, deliveries, nil
}

func (r *rabbitMQConn) Close() error {
	r.Lock()
	defer r.Unlock()

	select {
	case <-r.close:
		return nil
	default:
		close(r.close)
		r.connected = false
	}

	return r.Connection.Close()
}

func (r *rabbitMQConn) connect(secure bool, config *amqp.Config) error {
	// try connect
	if err := r.tryConnect(secure, config); err != nil {
		return err
	}

	// connected
	r.Lock()
	r.connected = true
	r.Unlock()

	// create reconnect loop
	go r.reconnect(secure, config)
	return nil
}

func (r *rabbitMQConn) reconnect(secure bool, config *amqp.Config) {
	// skip first connect
	var connect bool

	for {
		if connect {
			// try reconnect
			if err := r.tryConnect(secure, config); err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			// connected
			r.Lock()
			r.connected = true
			r.Unlock()
			//unblock resubscribe cycle - close channel
			//at this point channel is created and unclosed - close it without any additional checks
			close(r.waitConnection)
		}

		connect = true
		notifyClose := make(chan *amqp.Error)
		r.Connection.NotifyClose(notifyClose)

		// block until closed
		select {
		case <-notifyClose:
			// block all resubscribe attempt - they are useless because there is no connection to rabbitmq
			// create channel 'waitConnection' (at this point channel is nil or closed, create it without unnecessary checks)
			r.Lock()
			r.connected = false
			r.waitConnection = make(chan struct{})
			r.Unlock()
		case <-r.close:
			return
		}
	}
}

func (r *rabbitMQConn) tryConnect(secure bool, config *amqp.Config) error {
	var err error

	var url string

	if len(r.opts.Addr) > 0 && regexp.MustCompile("^amqp(s)?://.*").MatchString(r.opts.Addr[0]) {
		url = r.opts.Addr[0]
	} else {
		url = "amqp://guest:guest@127.0.0.1:5672"
	}

	if secure || config.TLSClientConfig != nil || strings.HasPrefix(url, "amqps://") {
		if config.TLSClientConfig == nil {
			config.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		url = strings.Replace(url, "amqp://", "amqps://", 1)
	}

	r.Connection, err = amqp.DialConfig(url, *config)

	if err != nil {
		return err
	}

	if r.Channel, err = newRabbitChannel(r.Connection, r.opts.RabbitMq.PrefetchCount, r.opts.RabbitMq.PrefetchGlobal); err != nil {
		return err
	}

	if r.opts.RabbitMq.DurableExchange {
		r.Channel.DeclareDurableExchange(r.opts.RabbitMq.ExchangeType, r.opts.RabbitMq.ExchangeKey)
	} else {
		r.Channel.DeclareExchange(r.opts.RabbitMq.ExchangeType, r.opts.RabbitMq.ExchangeKey)
	}
	r.ExchangeChannel, err = newRabbitChannel(r.Connection, r.opts.RabbitMq.PrefetchCount, r.opts.RabbitMq.PrefetchGlobal)

	return err
}