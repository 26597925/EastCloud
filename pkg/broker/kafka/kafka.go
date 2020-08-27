package kafka

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/prometheus/common/log"
	"sapi/pkg/broker"
	"sync"
)

type Producer struct {
	opts *broker.Options

	c sarama.Client
	p sarama.SyncProducer

	connected bool
}

func (pr *Producer) isConnected() bool {
	return pr.connected
}

func (pr *Producer) Connect() error {
	if pr.isConnected() {
		return nil
	}

	if pr.c != nil {
		pr.connected = true
		return nil
	}

	addr := pr.opts.Addr
	if len(addr) == 0 {
		addr = []string{"127.0.0.1:9092"}
	}

	config := pr.opts.Kafka
	if  config == nil {
		config = sarama.NewConfig()
	}
	if !config.Version.IsAtLeast(sarama.V0_11_0_0) {
		config.Version = sarama.V0_11_0_0
	}
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	c, err := sarama.NewClient(addr, config)
	if err != nil {
		return err
	}

	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		return err
	}

	pr.c = c
	pr.p = p
	pr.connected = true

	return nil
}

func (pr *Producer) Disconnect() error {
	if !pr.isConnected() {
		return nil
	}
	pr.p.Close()

	if err := pr.c.Close(); err != nil {
		return err
	}
	pr.connected = false
	return nil
}

func (pr *Producer) Publish(topic string, msg *broker.Message) error {
	if !pr.isConnected() {
		return errors.New("[kafka] broker not connected")
	}

	var headers []sarama.RecordHeader
	for k,v := range msg.Header {
		var header sarama.RecordHeader
		header.Key = []byte(k)
		header.Value = v.([]byte)
		headers = append(headers, header)
	}

	_, _, err := pr.p.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Headers: headers,
		Value: sarama.ByteEncoder(msg.Body),
	})

	return err
}

type Consumer struct {
	opts broker.SubscribeOptions

	cg   sarama.ConsumerGroup
	topic    string

	errorHandler broker.Handler
}

func (cr *Consumer) GroupID() string {
	return cr.opts.GroupID
}

func (cr *Consumer) Topic() string {
	return cr.topic
}

func (cr *Consumer) Subscribe(topic string, handler broker.Handler) error {
	cr.topic = topic
	h := &consumerGroupHandler{
		opts:    cr.opts,
		cg:      cr.cg,
		handler: handler,
		errorHandler: cr.errorHandler,
	}

	ctx := context.Background()
	topics := []string{topic}
	go func() {
		for {
			select {
			case err := <-cr.cg.Errors():
				if err != nil {
					if cr.errorHandler != nil {
						msg := &Event{
							topic: topic,
							err: err,
						}
						cr.errorHandler(msg)
					}

					log.Errorf("k_consumer errors:", err)
				}
			default:
				err := cr.cg.Consume(ctx, topics, h)
				switch err {
				case sarama.ErrClosedConsumerGroup:
					return
				case nil:
					continue
				default:
					if cr.errorHandler != nil {
						msg := &Event{
							topic: topic,
							err: err,
						}
						cr.errorHandler(msg)
					}

					log.Error(err)
				}
			}
		}
	}()

	return nil
}

func (cr *Consumer) Unsubscribe() error {
	return cr.cg.Close()
}

func NewBroker(opts *broker.Options) broker.Broker {
	return &kBroker{
		opts: opts,
	}
}

type kBroker struct {
	opts *broker.Options

	pr   *Producer
	crs  []*Consumer

	sc []sarama.Client
	scMutex   sync.RWMutex
}

func (k *kBroker) Options() *broker.Options {
	return k.opts
}

func (k *kBroker) NewProducer() (broker.Producer, error) {
	k.pr = &Producer{
		opts: k.opts,
		connected: false,
	}

	return k.pr, k.pr.Connect()
}

func (k *kBroker) NewConsumer(opts ...broker.SubscribeOption) (broker.Consumer, error) {
	opt := broker.SubscribeOptions{
		AutoAck: true,
		GroupID:   uuid.New().String(),
	}
	for _, o := range opts {
		o(&opt)
	}

	addr := k.opts.Addr
	if len(addr) == 0 {
		addr = []string{"127.0.0.1:9092"}
	}

	config := k.opts.Kafka
	if  config == nil {
		config = sarama.NewConfig()
	}
	if !config.Version.IsAtLeast(sarama.V0_11_0_0) {
		config.Version = sarama.V0_11_0_0
	}
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	cs, err := sarama.NewClient(addr, config)
	if err != nil {
		return nil, err
	}
	k.scMutex.Lock()
	defer k.scMutex.Unlock()
	k.sc = append(k.sc, cs)

	cg, err := sarama.NewConsumerGroupFromClient(opt.GroupID, cs)
	if err != nil {
		return nil, err
	}

	consumer := &Consumer{
		opts: opt,
		cg: cg,
		errorHandler: k.opts.ErrorHandler,
	}
	k.crs = append(k.crs, consumer)

	return consumer, nil
}

func (k *kBroker) Close() error {
	k.scMutex.Lock()
	defer k.scMutex.Unlock()

	for _,cr := range k.crs{
		cr.Unsubscribe()
	}

	for _, client := range k.sc {
		client.Close()
	}
	k.sc = nil

	return k.pr.Disconnect()
}

func (k *kBroker) Type() string {
	return "kafka"
}