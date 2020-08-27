package client

import (
	"context"
	"net/http"
)

var HTTPClient ContextKey

type ContextKey struct{}

func ContextClient(ctx context.Context) *http.Client {
	if ctx != nil {
		if hc, ok := ctx.Value(HTTPClient).(*http.Client); ok {
			return hc
		}
	}
	return http.DefaultClient
}

func NewClient(ctx context.Context, src Source) *http.Client {
	if src == nil {
		return ContextClient(ctx)
	}
	return &http.Client{
		Transport: &Transport{
			Base:   ContextClient(ctx).Transport,
			Source: ReuseTokenSource(nil, src),
		},
	}
}

