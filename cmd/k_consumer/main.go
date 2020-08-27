package main

import (
	"fmt"
	"sapi/pkg/broker"
	"sapi/pkg/broker/kafka"
)

func main() {
	kafka := kafka.NewBroker(&broker.Options{})
	csr, err := kafka.NewConsumer()
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
