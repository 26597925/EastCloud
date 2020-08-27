package proxy

type Proxy struct {
	opt *Options
}

func NewProxy(opt *Options) *Proxy {
	p := &Proxy{
		opt: opt,
	}

	return p
}

func (p *Proxy) Start() {

}

func (p *Proxy) Stop() {

}