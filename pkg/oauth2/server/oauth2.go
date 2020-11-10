package server

import (
	"context"
	"fmt"
	"github.com/26597925/EastCloud/pkg/oauth2/api"
	"github.com/26597925/EastCloud/pkg/oauth2/server/errors"
	"github.com/26597925/EastCloud/pkg/oauth2/server/generates"
	"github.com/26597925/EastCloud/pkg/oauth2/server/stores"
	"github.com/26597925/EastCloud/pkg/oauth2/server/token"
	"github.com/26597925/EastCloud/pkg/util/crypto"
	"net/http"
	"net/url"
	"strings"
)

type Oauth2 struct {
	opts *Options
	gen *Generates

	ClientInfoHandler            ClientInfoHandler
	UserAuthorizationHandler     UserAuthorizationHandler
	PasswordAuthorizationHandler PasswordAuthorizationHandler
	AuthorizeScopeHandler        AuthorizeScopeHandler
	ValidateClientHandler		 ValidateClientHandler
}

func NewOauth2(opts *Options) *Oauth2{
	if opts.tokenConfigs == nil {
		opts.tokenConfigs = GetDefaultTokenConfig()
	}

	cs := &stores.ClientStore{
		Store: GetDefaultStore(),
	}

	ts :=&stores.TokenStore{
		Store: GetDefaultStore(),
	}

	gen    := &Generates{
		codeExp: opts.codeExp,
		TokenType: opts.TokenType,
		MacAlgorithm: opts.MacAlgorithm,
		AllowedGrantTypes: opts.AllowedGrantTypes,
		tokenConfigs:	opts.tokenConfigs,
		refreshConfig:  opts.refreshConfig,
		ClientStore: 	cs,
		TokenStore:  	ts,
		AuthorizeGenerate: generates.NewAuthorize(),
		AccessGenerate: generates.NewBaseAccess(),
		ValidateURIHandler: DefValidateURIHandler,
	}

	oauth2 := &Oauth2{
		opts:	opts,
		gen:	gen,
	}

	oauth2.ClientInfoHandler = ClientBasicHandler
	oauth2.UserAuthorizationHandler = DefUserAuthorizationHandler
	oauth2.PasswordAuthorizationHandler = DefPasswordAuthorizationHandler

	return oauth2
}

func (o *Oauth2) SetClientStore(clientStore *stores.ClientStore) {
	o.gen.ClientStore = clientStore
}

func (o *Oauth2) SetTokenStore(tokenStore *stores.TokenStore) {
	o.gen.TokenStore = tokenStore
}

func (o *Oauth2) SetAuthorizeGenerate(authorizeGenerate generates.AuthorizeGenerate) {
	o.gen.AuthorizeGenerate = authorizeGenerate
}

func (o *Oauth2) SetAccessGenerate(accessGenerate generates.AccessGenerate) {
	o.gen.AccessGenerate = accessGenerate
}

func (o *Oauth2) SetValidateURIHandler(validateURIHandler ValidateURIHandler) {
	o.gen.ValidateURIHandler = validateURIHandler
}

func (o *Oauth2) SetExtensionFieldsHandler(extensionFieldsHandler ExtensionFieldsHandler) {
	o.gen.ExtensionFieldsHandler = extensionFieldsHandler
}

func (o *Oauth2) SetClientAuthorizedHandler(clientAuthorizedHandler ClientAuthorizedHandler) {
	o.gen.ClientAuthorizedHandler = clientAuthorizedHandler
}

func (o *Oauth2) SetClientScopeHandler(clientScopeHandler ClientScopeHandler) {
	o.gen.ClientScopeHandler = clientScopeHandler
}

func (o *Oauth2) SetRefreshingScopeHandler(refreshingScopeHandler RefreshingScopeHandler) {
	o.gen.RefreshingScopeHandler = refreshingScopeHandler
}

func (o *Oauth2) SaveClient(ctx context.Context, client *stores.Client) {
	o.gen.ClientStore.Set(ctx, client.ID, client)
}

func (o *Oauth2) HandleAuthorizeRequest(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &Request{
		Request:r,
		AllowedResponseTypes:         o.opts.AllowedResponseTypes,
		AllowGetAccessRequest:        o.opts.AllowGetAccessRequest,
		ClientInfoHandler:            o.ClientInfoHandler,
		PasswordAuthorizationHandler: o.PasswordAuthorizationHandler,
	}
	res := &Response{
		Writer: w,
	}

	ar, err := req.ValidationAuthorizeRequest()
	if err != nil {
		return o.redirectError(res, ar, err)
	}

	userID, err := o.UserAuthorizationHandler(w, r)
	if err != nil {
		return o.redirectError(res, ar, err)
	} else if userID == "" {
		return nil
	}
	ar.UserID = userID

	if fn := o.AuthorizeScopeHandler; fn != nil {
		scope, err := fn(w, r)
		if err != nil {
			return err
		} else if scope != "" {
			ar.Scope = scope
		}
	}

	tk, err := o.gen.GetAuthorizeToken(ctx, ar)
	if err != nil {
		return o.redirectError(res, ar, err)
	}

	if ar.RedirectURI == "" {
		client, err := o.gen.ClientStore.GetClient(ctx, ar.ClientID)
		if err != nil {
			return err
		}
		ar.RedirectURI = client.Domain
	}

	uri, err := o.getRedirectURI(ar, tk.GetAuthorizeData())
	if err != nil {
		return err
	}

	return res.Redirect(uri)
}

func (o *Oauth2) HandleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	req := &Request{
		Request:r,
		AllowedResponseTypes:         o.opts.AllowedResponseTypes,
		AllowGetAccessRequest:        o.opts.AllowGetAccessRequest,
		ClientInfoHandler:            o.ClientInfoHandler,
		PasswordAuthorizationHandler: o.PasswordAuthorizationHandler,
	}

	res := &Response{
		Writer: w,
	}

	tr, err := req.ValidationTokenRequest()
	if err != nil {
		return res.OutputTokenError(err)
	}

	t, err := o.gen.GetAccessToken(ctx, tr)
	if err != nil {
		return res.OutputTokenError(err)
	}

	return res.OutputToken(t.BuildToken())
}

// ValidationBearerToken validation the bearer tokens
// https://tools.ietf.org/html/rfc6750
func (o *Oauth2) ValidationBearerToken(r *http.Request) (*token.Token, error) {
	ctx := r.Context()
	req := &Request{
		Request:r,
	}

	accessToken, ok := req.BearerAuth()
	if !ok {
		return nil, errors.ErrInvalidAccessToken
	}

	return o.gen.LoadAccessToken(ctx, accessToken)
}

func (o *Oauth2) ValidationMACToken(r *http.Request) (*token.Token, error){
	ctx := r.Context()
	req := &Request{
		Request:r,
	}
	mac, err := req.MacAuth()
	if err != nil {
		return nil, errors.ErrInvalidAccessToken
	}
	if mac.id == "" || mac.mac == "" || mac.nonce == "" {
		return nil, errors.ErrInvalidAccessToken
	}
	if mac.ts <= 0 {
		return nil, errors.ErrExpiredAccessToken
	}

	hostName := r.Host
	port := "80"
	len := strings.IndexByte(r.Host, ':')
	if len > -1 {
		hostName = r.Host[:len]
		port = r.Host[len+1:]
	}

	if fn := o.ValidateClientHandler; fn != nil {
		allowed, err := fn(mac.nonce, mac.ts, mac.ext)
		if err != nil {
			return nil, err
		} else if !allowed {
			return nil, errors.ErrInvalidRequest
		}
	}

	token, err := o.gen.LoadAccessToken(ctx, mac.id)
	if err != nil {
		return nil, errors.ErrInvalidAccessToken
	}

	text := fmt.Sprintf("%d\n%s\n%s\n%s\n%s\n%s\n%s\n", mac.ts, mac.nonce, r.Method, r.RequestURI, hostName, port, mac.ext)
	if strings.EqualFold(token.MacAlgorithm, "hmac-sha-1") {
		macText := crypto.HmacSha1(token.MacKey, text)
		if macText != mac.mac {
			return nil, errors.ErrInvalidAccessToken
		}
	}
	if strings.EqualFold(token.MacAlgorithm, "hmac-sha-256") {
		macText := crypto.HmacSha256(token.MacKey, text)
		if macText != mac.mac {
			return nil, errors.ErrInvalidAccessToken
		}
	}

	return token, nil
}

func (o *Oauth2) getRedirectURI(req *AuthorizeRequest, data map[string]interface{}) (string, error) {
	u, err := url.Parse(req.RedirectURI)
	if err != nil {
		return "", err
	}

	q := u.Query()
	if req.State != "" {
		q.Set("state", req.State)
	}

	for k, v := range data {
		q.Set(k, fmt.Sprint(v))
	}

	switch req.ResponseType {
	case api.Code:
		u.RawQuery = q.Encode()
	case api.Token:
		u.RawQuery = ""
		fragment, err := url.QueryUnescape(q.Encode())
		if err != nil {
			return "", err
		}
		u.Fragment = fragment
	}

	return u.String(), nil
}

func (o *Oauth2) redirectError(res *Response, req *AuthorizeRequest, err error) error {
	if req == nil {
		return err
	}
	data, _:= res.ErrorData(err)
	uri, err := o.getRedirectURI(req, data)
	if err != nil {
		return err
	}

	return res.Redirect(uri)
}