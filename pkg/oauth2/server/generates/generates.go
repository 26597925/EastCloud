package generates

import (
	"context"
	"sapi/pkg/oauth2/server/token"
)

type (
	AuthorizeGenerate interface {
		Token(ctx context.Context, token * token.Token) (code string, err error)
	}

	AccessGenerate interface {
		Token(ctx context.Context, token * token.Token, isGenRefresh bool) (access, refresh string, err error)
	}
)