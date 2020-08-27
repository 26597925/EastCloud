package token

import (
	"sapi/pkg/oauth2/api"
	"strings"
	"time"
)

type Token struct {
	TokenType		 string
	GrantType 		 api.GrantType
	ResponseType 	 api.ResponseType
	ClientID         string
	UserID           string
	RedirectURI      string
	Scope            string
	Code             string
	CodeCreateAt     time.Time
	CodeExpiresIn    time.Duration
	Access           string
	AccessCreateAt   time.Time
	AccessExpiresIn  time.Duration
	Refresh          string
	RefreshCreateAt  time.Time
	RefreshExpiresIn time.Duration

	MacKey	string
	MacAlgorithm	string
	FieldsValue 	 map[string]interface{}
}

func (t *Token) GetAuthorizeData() map[string]interface{} {
	if t.ResponseType == api.Code {
		return map[string]interface{}{
			"code": t.Code,
		}
	}
	return t.BuildToken()
}

func (t *Token) BuildToken() map[string]interface{} {
	data := map[string]interface{}{
		"access_token": t.Access,
		"token_type":   t.TokenType,
		"expires_in":   int64(t.AccessExpiresIn / time.Second),
	}

	if strings.EqualFold(t.TokenType, "mac") {
		data["mac_key"] = t.MacKey
		data["mac_algorithm"] = t.MacAlgorithm
	}

	if scope := t.Scope; scope != "" {
		data["scope"] = scope
	}

	if refresh := t.Refresh; refresh != "" {
		data["refresh_token"] = refresh
	}

	for k, v := range t.FieldsValue {
		if _, ok := data[k]; ok {
			continue
		}
		data[k] = v
	}

	return data
}