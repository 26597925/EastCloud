package registry

type Result struct {
	Action  string
	Service *Service
}

type Watcher interface {
	// Next is a blocking call
	Next() (*Result, error)
	Stop()
}
