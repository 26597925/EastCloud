package client

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"sapi/pkg/oauth2/api"
	"strings"
)

type Oauth2 struct {
	opts *Options
}

func NewOauth2(opts *Options) *Oauth2{
	return &Oauth2{opts:opts}
}

func (o *Oauth2) AuthCodeURL(state string, opts ...AuthCodeOption) string{
	var buf bytes.Buffer
	buf.WriteString(o.opts.Endpoint.AuthURL)
	v := url.Values{
		"response_type": {api.Code.String()},
		"client_id":     {o.opts.ClientID},
	}
	if o.opts.RedirectURL != "" {
		v.Set("redirect_uri", o.opts.RedirectURL)
	}
	if len(o.opts.Scopes) > 0 {
		v.Set("scope", strings.Join(o.opts.Scopes, " "))
	}
	if state != "" {
		v.Set("state", state)
	}
	for _, opt := range opts {
		opt.setValue(v)
	}
	if strings.Contains(o.opts.Endpoint.AuthURL,"?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(v.Encode())
	return buf.String()
}

func (o *Oauth2) AuthorizationCodeToken(ctx context.Context, code string, opts ...AuthCodeOption) (*Token, error) {
	v := url.Values{
		"grant_type": {api.AuthorizationCode.String()},
		"code":       {code},
	}
	if o.opts.RedirectURL != "" {
		v.Set("redirect_uri", o.opts.RedirectURL)
	}
	for _, opt := range opts {
		opt.setValue(v)
	}
	return RequestToken(ctx, o.opts, v)
}

func (o *Oauth2) PasswordCredentialsToken(ctx context.Context, username, password string) (*Token, error) {
	v := url.Values{
		"grant_type": {api.PasswordCredentials.String()},
		"username":   {username},
		"password":   {password},
	}
	if len(o.opts.Scopes) > 0 {
		v.Set("scope", strings.Join(o.opts.Scopes, " "))
	}
	return  RequestToken(ctx, o.opts, v)
}

func (o *Oauth2) ClientCredentialsToken(ctx context.Context, opts ...AuthCodeOption) (*Token, error)  {
	v := url.Values{
		"grant_type": {api.ClientCredentials.String()},
	}
	if len(o.opts.Scopes) > 0 {
		v.Set("scope", strings.Join(o.opts.Scopes, " "))
	}

	for _, opt := range opts {
		opt.setValue(v)
	}

	return RequestToken(ctx, o.opts, v)
}

func (o *Oauth2) TokenSource(ctx context.Context, t *Token) Source {
	return NewTokenSource(ctx, o.opts, t)
}

func (o *Oauth2) Client(ctx context.Context, t *Token) *http.Client {
	return NewClient(ctx, o.TokenSource(ctx, t))
}

func (o *Oauth2) StaticClient(ctx context.Context, accessToken string) *http.Client {
	return NewClient(ctx, StaticTokenSource(&Token{
		AccessToken: accessToken,
	}))
}