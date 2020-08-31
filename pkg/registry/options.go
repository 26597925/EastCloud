package registry

type Options struct {
	Prefix    string
	TTL 	  int64
	Timeout   int64
}

type Option func(*Options)