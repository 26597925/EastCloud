package rabbitmq

import (
	"errors"
	"github.com/26597925/EastCloud/pkg/broker"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

const (
	TOPIC  = "topic"
	DELAYED = "x-delayed-message"
	)

type rBroker struct {
	opts *broker.Options
	conn   *rabbitMQConn

	mtx            sync.Mutex
	wg             sync.WaitGroup
}

type Consumer struct {
	r            *rBroker
	mtx          sync.Mutex
	mayRun       bool
	opts         broker.SubscribeOptions
	topic        string
	ch           *rabbitMQChannel
	durableQueue bool
	queueArgs    map[string]interface{}
	fn           func(msg amqp.Delivery)
	headers      map[string]interface{}
}

type Producer struct {
	conn  *rabbitMQConn
	publishing amqp.Publishing
}

func NewBroker(opts *broker.Options) broker.Broker {
	return &rBroker{
		opts: opts,
	}
}

func (r *rBroker) Options() *broker.Options {
	return r.opts
}

func (r *rBroker) NewProducer() (broker.Producer, error) {
	err := r.connect()

	producer := &Producer{conn:r.conn}
	return producer, err
}

func (r *rBroker) NewConsumer(opts ...broker.SubscribeOption) (broker.Consumer, error) {
	err := r.connect()
	opt := broker.SubscribeOptions{
		GroupID: uuid.New().String(),
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	csr := &Consumer{opts: opt, mayRun: true, r: r}

	return csr, err
}

func (r *rBroker) connect() error {
	if r.conn == nil {
		r.conn = newRabbitMQConn(r.opts)
	}

	conf := amqp.Config{
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
	}

	if r.opts.RabbitMq == nil {
		r.opts.RabbitMq = &broker.RabbitMq{
			ExchangeType: "topic",
			ExchangeKey:  "Spi",
		}
	}

	if r.opts.RabbitMq.ExchangeKey == "" {
		r.opts.RabbitMq.ExchangeKey = "Spi"
	}

	if r.opts.RabbitMq.Authentication != nil {
		conf.SASL = []amqp.Authentication{r.opts.RabbitMq.Authentication}
	}

	conf.TLSClientConfig = r.opts.RabbitMq.TLSConfig

	return r.conn.Connect(r.opts.RabbitMq.Secure, &conf)
}

func (r *rBroker) Close() error {
	if r.conn == nil {
		return errors.New("connection is nil")
	}
	ret := r.conn.Close()
	r.wg.Wait() // wait all goroutines
	return ret
}

func (r *rBroker) Type() string {
	return "rabbitMq"
}

func (pr *Producer) SetPublish(publishing amqp.Publishing) {
	pr.publishing = publishing
}

func (pr *Producer) Publish(topic string, msg *broker.Message) error {
	pr.publishing.Body = msg.Body
	if pr.publishing.Headers == nil {
		pr.publishing.Headers = amqp.Table{}
	}

	for k, v := range msg.Header {
		pr.publishing.Headers[k] = v
	}

	return pr.conn.Publish(pr.conn.opts.RabbitMq.ExchangeKey, topic, pr.publishing)
}

func (cr *Consumer) GroupID() string {
	return cr.opts.GroupID
}

func (cr *Consumer) Topic() string {
	return cr.topic
}

func (cr *Consumer) SetDurableQueue(durableQueue bool) {
	cr.durableQueue = durableQueue
}

func (cr *Consumer)  SetQueueArgs (queueArgs map[string]interface{}) {
	cr.queueArgs = queueArgs
}

func (cr *Consumer) SetHeaders (headers map[string]interface{}) {
	cr.headers = headers
}

func (cr *Consumer) Subscribe(topic string, handler broker.Handler)  error {
	cr.topic = topic
	cr.fn = func(msg amqp.Delivery) {
		header := make(map[string]interface{})
		for k, v := range msg.Headers {
			header[k] = v
		}
		m := &broker.Message{
			Header: header,
			Body:   msg.Body,
		}
		p := &Event{d: msg, m: m, t: msg.RoutingKey}
		p.err = handler(p)
		if p.err == nil && !cr.opts.AutoAck {
			msg.Ack(false)
		} else if p.err != nil && !cr.opts.AutoAck {
			msg.Nack(false, true)
		}
	}

	go cr.resubscribe()

	return nil
}

func (cr *Consumer) Unsubscribe() error {
	cr.mtx.Lock()
	defer cr.mtx.Unlock()
	cr.mayRun = false
	if cr.ch != nil {
		return cr.ch.Close()
	}
	return nil
}

func (cr *Consumer) resubscribe() {
	minResubscribeDelay := 100 * time.Millisecond
	maxResubscribeDelay := 30 * time.Second
	expFactor := time.Duration(2)
	reSubscribeDelay := minResubscribeDelay
	//loop until unsubscribe
	for {
		cr.mtx.Lock()
		mayRun := cr.mayRun
		cr.mtx.Unlock()
		if !mayRun {
			// we are unsubscribed, showdown routine
			return
		}

		select {
		//check shutdown case
		case <-cr.r.conn.close:
			//yep, its shutdown case
			return
			//wait until we reconect to rabbit
		case <-cr.r.conn.waitConnection:
		}

		// it may crash (panic) in case of Consume without connection, so recheck it
		cr.r.mtx.Lock()
		if !cr.r.conn.connected {
			cr.r.mtx.Unlock()
			continue
		}

		ch, sub, err := cr.r.conn.Consume(
			cr.opts.GroupID,
			cr.topic,
			cr.headers,
			cr.queueArgs,
			cr.opts.AutoAck,
			cr.durableQueue,
		)

		cr.r.mtx.Unlock()
		switch err {
		case nil:
			reSubscribeDelay = minResubscribeDelay
			cr.mtx.Lock()
			cr.ch = ch
			cr.mtx.Unlock()
		default:
			if reSubscribeDelay > maxResubscribeDelay {
				reSubscribeDelay = maxResubscribeDelay
			}
			time.Sleep(reSubscribeDelay)
			reSubscribeDelay *= expFactor
			continue
		}
		for d := range sub {
			cr.r.wg.Add(1)
			cr.fn(d)
			cr.r.wg.Done()
		}
	}
}