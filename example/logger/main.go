package main

import (
	"fmt"
	"sapi/pkg/logger"
	"time"
)

type user struct {
	name string
	age int
}

func main () {
	field := make(map[string]interface{})

	var value time.Duration = 3
	field["cost"] = fmt.Sprintf("%.3f", float64(value.Round(time.Microsecond))/float64(time.Millisecond))

	user := &user{
		name: "aaa",
		age:  10,
	}

	logger.Info(user)

}
