package proxy

import "net/http"

func (p *Proxy) newHTTPServer() *http.Server{
	return &http.Server{
			Addr: p.opt.Addr,
		}
}

func (p *Proxy) startHTTP() {
	s := p.newHTTPServer()
	err := s.ListenAndServe()
	if err != nil {
		//log.Fatalf("start http listeners failed with %+v", err)
	}


}
