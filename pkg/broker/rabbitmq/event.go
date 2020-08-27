package rabbitmq

import (
	"github.com/streadway/amqp"
	"sapi/pkg/broker"
)

type Event struct {
	d   amqp.Delivery
	m   *broker.Message
	t   string
	err error
}

func (evt *Event) Ack() error {
	return evt.d.Ack(false)
}

func (evt *Event) Error() error {
	return evt.err
}

func (evt *Event) Topic() string {
	return evt.t
}

func (evt *Event) Message() *broker.Message {
	return evt.m
}


