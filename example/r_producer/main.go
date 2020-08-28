package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"sapi/pkg/broker"
	"sapi/pkg/broker/rabbitmq"
)

//https://blog.csdn.net/qq_26656329/article/details/77891154
func main () {
	rabbitMq := rabbitmq.NewBroker(&broker.Options{RabbitMq:&broker.RabbitMq{ExchangeType:"x-delayed-message"}})
	producer, error := rabbitMq.NewProducer()

	pro := producer.(*rabbitmq.Producer)
	pro.SetPublish(amqp.Publishing{
		Headers: amqp.Table{
			"x-delay": 8000,
		},
	})

	error = producer.Publish("1111", &broker.Message{
		Header: map[string]interface{}{"aa":"bb"},
		Body: []byte("asdasd"),
	})

	fmt.Println(error)

}
