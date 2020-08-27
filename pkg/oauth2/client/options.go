package client

import (
	"net/url"
)

var (
	AccessTypeOnline  AuthCodeOption = SetAuthURLParam("access_type", "online")
	AccessTypeOffline AuthCodeOption = SetAuthURLParam("access_type", "offline")
	ApprovalForce AuthCodeOption = SetAuthURLParam("prompt", "consent")
)

type AuthStyle int
const (
	AuthStyleInParams AuthStyle = 0
	AuthStyleInHeader AuthStyle = 1
)

type Endpoint struct {
	AuthURL  string
	TokenURL string
}

type AuthCodeOption interface {
	setValue(url.Values)
}

type setParam struct{ k, v string }

func (p setParam) setValue(m url.Values) { m.Set(p.k, p.v) }

func SetAuthURLParam(key, value string) AuthCodeOption {
	return setParam{key, value}
}

type Options struct {
	ClientID     string
	ClientSecret string
	Domain string

	Endpoint Endpoint
	RedirectURL string
	Scopes []string
	AuthStyle AuthStyle
}