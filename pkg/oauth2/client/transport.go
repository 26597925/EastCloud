package client

import (
	"errors"
	"net/http"
	"sapi/pkg/logger"
	"sync"
)
var cancelOnce sync.Once

type Transport struct {
	Source Source
	Base   http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBodyClosed := false
	if req.Body != nil {
		defer func() {
			if !reqBodyClosed {
				req.Body.Close()
			}
		}()
	}

	if t.Source == nil {
		return nil, errors.New("oauth2_client: Transport's Source is nil")
	}
	token, err := t.Source.Token()
	if err != nil {
		return nil, err
	}

	req2 := cloneRequest(req)
	token.SetAuthHeader(req2)

	reqBodyClosed = true
	return t.base().RoundTrip(req2)
}

func (t *Transport) CancelRequest(req *http.Request) {
	cancelOnce.Do(func() {
		logger.Info("deprecated: golang.org/x/oauth2_client: Transport.CancelRequest no longer does anything; use contexts")
	})
}

func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

func cloneRequest(r *http.Request) *http.Request {
	r2 := new(http.Request)
	*r2 = *r
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
