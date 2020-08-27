package registry

var (
	DefaultPrefix = "/sapi/registry"
)

type Options struct {
	Endpoints []string
	Timeout   int
	TTL       int
	Prefix    string

	BasicAuth bool
	Username  string
	Password  string

	CertFile  string
	KeyFile   string
	CaCert    string
}

type Option func(*Options)