package server

import (
	"github.com/26597925/EastCloud/pkg/oauth2/api"
	"net/url"
	"time"
)

var (
	AccessTypeOnline  AuthCodeOption = SetAuthURLParam("access_type", "online")
	AccessTypeOffline AuthCodeOption = SetAuthURLParam("access_type", "offline")
	ApprovalForce AuthCodeOption = SetAuthURLParam("prompt", "consent")
)

var (
	DefaultCodeExp               = time.Minute * 10
	DefaultAuthorizeCodeTokenCfg = &TokenConfig{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 3, IsGenerateRefresh: true}
	DefaultImplicitTokenCfg      = &TokenConfig{AccessTokenExp: time.Hour * 1}
	DefaultPasswordTokenCfg      = &TokenConfig{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 7, IsGenerateRefresh: true}
	DefaultClientTokenCfg        = &TokenConfig{AccessTokenExp: time.Hour * 2}
	DefaultRefreshTokenCfg       = &RefreshingConfig{IsGenerateRefresh: true, IsRemoveAccess: true, IsRemoveRefreshing: true}
)

type AuthCodeOption interface {
	setValue(url.Values)
}

type setParam struct{ k, v string }

func (p setParam) setValue(m url.Values) { m.Set(p.k, p.v) }

func SetAuthURLParam(key, value string) AuthCodeOption {
	return setParam{key, value}
}

type TokenConfig struct {
	AccessTokenExp time.Duration
	RefreshTokenExp time.Duration
	IsGenerateRefresh bool
}

type RefreshingConfig struct {
	AccessTokenExp time.Duration
	RefreshTokenExp time.Duration
	IsGenerateRefresh bool
	IsResetRefreshTime bool
	IsRemoveAccess bool
	IsRemoveRefreshing bool
}

type Options struct {
	TokenType             string
	AllowGetAccessRequest bool
	AllowedResponseTypes  []api.ResponseType
	AllowedGrantTypes     []api.GrantType

	codeExp           time.Duration
	tokenConfigs	map[api.GrantType]*TokenConfig
	refreshConfig   *RefreshingConfig

	MacAlgorithm	string
}

func GetDefaultTokenConfig() map[api.GrantType]*TokenConfig {
	tokenConfigs := make(map[api.GrantType]*TokenConfig)
	tokenConfigs[api.AuthorizationCode] = DefaultAuthorizeCodeTokenCfg
	tokenConfigs[api.Implicit] = DefaultImplicitTokenCfg
	tokenConfigs[api.PasswordCredentials] = DefaultPasswordTokenCfg
	tokenConfigs[api.ClientCredentials] = DefaultClientTokenCfg

	return tokenConfigs
}

func (opt *Options) SetTokenConfig(gt api.GrantType, tc *TokenConfig) {
	opt.tokenConfigs[gt] = tc
}

func (opt *Options) SetRefreshConfig(refreshingConfig *RefreshingConfig) {
	opt.refreshConfig = refreshingConfig
}

func (opt *Options) SetCodeExp(codeExp time.Duration) {
	opt.codeExp = codeExp
}