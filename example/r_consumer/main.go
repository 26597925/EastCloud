package main

import (
	"fmt"
	"sapi/pkg/broker"
	"sapi/pkg/broker/rabbitmq"
)

func main () {
	rabbitMq := rabbitmq.NewBroker(&broker.Options{RabbitMq:&broker.RabbitMq{ExchangeType:rabbitmq.DELAYED}})
	csr, err := rabbitMq.NewConsumer()
	fmt.Println(err)

	csr.Subscribe("1111", func(evt broker.Event) error {
		fmt.Println(evt.Message().Header)
		fmt.Println(string(evt.Message().Body))
		return nil
	})

	for true {
		select {

		}
	}
}
