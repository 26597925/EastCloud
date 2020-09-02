package registry

type Result struct {
	Type  	int
	Key   	string
	Service *Service
}

type Watcher interface {
	// Next is a blocking call
	Next() (*Result, error)
	Stop()
}
