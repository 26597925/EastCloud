package registry

type Result struct {
	Type  	int
	Service *Service
}

type Watcher interface {
	// Next is a blocking call
	Next() (*Result, error)
	Stop()
}
