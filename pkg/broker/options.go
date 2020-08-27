package broker

import (
	"crypto/tls"
	"github.com/Shopify/sarama"
	"github.com/streadway/amqp"
)

type Options struct {
	Addr []string
	Kafka *sarama.Config
	RabbitMq *RabbitMq
	ErrorHandler Handler
}

type Option func(*Options)

type SubscribeOptions struct {
	AutoAck bool
	GroupID string
}

type SubscribeOption func(*SubscribeOptions)

type RabbitMq struct {
	ExchangeType 	string
	ExchangeKey 	string
	DurableExchange bool
	PrefetchCount 	int
	PrefetchGlobal 	bool

	Secure 	  bool
	TLSConfig *tls.Config
	Authentication amqp.Authentication
}

//https://github.com/gaopengfei123123/go_study/blob/master/src/mq_demo/mqhandler/rabbitmq/server.go