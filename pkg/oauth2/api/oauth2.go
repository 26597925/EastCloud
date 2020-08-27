package api

import "context"

type GrantType string
const (
	AuthorizationCode   GrantType = "authorization_code"
	PasswordCredentials GrantType = "password"
	ClientCredentials   GrantType = "client_credentials"
	Refreshing          GrantType = "refresh_token"
	Implicit            GrantType = "__implicit"
)

func (gt GrantType) String() string {
	if gt == AuthorizationCode ||
		gt == PasswordCredentials ||
		gt == ClientCredentials ||
		gt == Refreshing {
		return string(gt)
	}
	return ""
}

type ResponseType string
const (
	Code  ResponseType = "code"
	Token ResponseType = "token"
)

func (rt ResponseType) String() string {
	return string(rt)
}

type Generates interface {
	Token(ctx context.Context, isGenRefresh bool) (string, string, error)
}