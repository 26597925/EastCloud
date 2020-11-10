package kafka

import (
	"github.com/26597925/EastCloud/pkg/broker"
	"github.com/Shopify/sarama"
	"github.com/prometheus/common/log"
)

type Event  struct {
	topic    string
	message *broker.Message
	err 	 error

	cm   *sarama.ConsumerMessage
	sess    sarama.ConsumerGroupSession
}

func (evt *Event) Topic() string {
	return evt.topic
}

func (evt *Event) Message() *broker.Message {
	return evt.message
}

func (evt *Event) Error() error {
	return evt.err
}

func (evt *Event) Ack() error {
	evt.sess.MarkMessage(evt.cm, "")
	return nil
}

type consumerGroupHandler struct {
	opts broker.SubscribeOptions
	cg      sarama.ConsumerGroup
	sess    sarama.ConsumerGroupSession

	handler broker.Handler
	errorHandler broker.Handler
}

func (*consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		headers := make(map[string]interface{})
		for _,v := range msg.Headers {
			headers[string(v.Key)] = string(v.Value)
		}

		event := &Event{
			topic: msg.Topic,
			message: &broker.Message{
				Header: headers,
				Body:   msg.Value,
			},
			cm: msg,
			sess: sess,
		}

		err := h.handler(event)
		if err == nil && h.opts.AutoAck {
			sess.MarkMessage(msg, "")
		} else if err != nil {
			event.err = err
			if h.errorHandler != nil {
				h.errorHandler(event)
			}
			log.Errorf("[kafka]: subscriber errors: %v", err)
		}
	}
	return nil
}