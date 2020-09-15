package scheduler

import (
	"errors"
	"sapi/pkg/util/rand"
)

const (
	RAND = iota
)

func SelectClient(scheduler int, clients []string) (string, error) {
	if scheduler == RAND {
		return RandSelector(clients)
	}

	return "", errors.New("no client")
}

func RandSelector(clients []string) (string, error) {
	l := len(clients)
	if 0 >= l {
		return "", errors.New("no client")
	}

	client := clients[rand.Uint32n(uint32(l))]
	return client, nil
}

