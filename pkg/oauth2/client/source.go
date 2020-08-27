package client

import (
	"context"
	"errors"
	"net/url"
	"sync"
)

type Source interface {
	Token() (*Token, error)
}

type staticTokenSource struct {
	t *Token
}

func StaticTokenSource(t *Token) Source {
	return staticTokenSource{t}
}

func (s staticTokenSource) Token() (*Token, error) {
	return s.t, nil
}

type tokenRefresher struct {
	ctx          context.Context // used to get HTTP requests
	opt  		 *Options
	refreshToken string
}

type reuseTokenSource struct {
	new Source

	mu sync.Mutex
	t  *Token
}

func ReuseTokenSource(t *Token, src Source) Source {
	if rt, ok := src.(*reuseTokenSource); ok {
		if t == nil {
			return rt
		}
		src = rt.new
	}
	return &reuseTokenSource{
		t:   t,
		new: src,
	}
}

func (s *reuseTokenSource) Token() (*Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.t.Valid() {
		return s.t, nil
	}
	t, err := s.new.Token()
	if err != nil {
		return nil, err
	}
	s.t = t
	return t, nil
}

func (tf *tokenRefresher) Token() (*Token, error) {
	if tf.refreshToken == "" {
		return nil, errors.New("oauth2_client: token expired and refresh token is not set")
	}

	tk, err := RequestToken(tf.ctx, tf.opt, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {tf.refreshToken},
	})

	if err != nil {
		return nil, err
	}
	if tf.refreshToken != tk.RefreshToken {
		tf.refreshToken = tk.RefreshToken
	}
	return tk, err
}

func NewTokenSource(ctx context.Context, opt  *Options, t *Token) Source{
	tkr := &tokenRefresher{
		ctx:  ctx,
		opt: opt,
	}
	if t != nil {
		tkr.refreshToken = t.RefreshToken
	}
	return &reuseTokenSource{
		t:   t,
		new: tkr,
	}
}
