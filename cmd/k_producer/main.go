package main

import (
	"fmt"
	"sapi/pkg/broker"
	"sapi/pkg/broker/kafka"
)

func main()  {
	kafka := kafka.NewBroker(&broker.Options{})
	producer, error := kafka.NewProducer()
	error = producer.Publish("1111", &broker.Message{
		Header: map[string]string{"aa":"bb"},
		Body: []byte("asdasd"),
	})
	error = kafka.Close()

	fmt.Println(error)
}
