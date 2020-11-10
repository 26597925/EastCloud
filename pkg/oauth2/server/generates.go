package server

import (
	"context"
	"github.com/26597925/EastCloud/pkg/logger"
	"github.com/26597925/EastCloud/pkg/oauth2/api"
	"github.com/26597925/EastCloud/pkg/oauth2/server/errors"
	"github.com/26597925/EastCloud/pkg/oauth2/server/generates"
	"github.com/26597925/EastCloud/pkg/oauth2/server/stores"
	"github.com/26597925/EastCloud/pkg/oauth2/server/token"
	"github.com/26597925/EastCloud/pkg/store"
	"github.com/26597925/EastCloud/pkg/store/buntdb"
	"strings"
	"time"
)

type Generates struct {
	codeExp          	  time.Duration
	TokenType       	  string
	MacAlgorithm		  string
	AllowedGrantTypes     []api.GrantType
	tokenConfigs		  map[api.GrantType]*TokenConfig
	refreshConfig   	  *RefreshingConfig

	ClientStore 		*stores.ClientStore
	TokenStore  		*stores.TokenStore

	AuthorizeGenerate generates.AuthorizeGenerate
	AccessGenerate	  generates.AccessGenerate

	ValidateURIHandler		ValidateURIHandler
	ExtensionFieldsHandler  ExtensionFieldsHandler

	ClientAuthorizedHandler      ClientAuthorizedHandler
	ClientScopeHandler			 ClientScopeHandler
	RefreshingScopeHandler 		 RefreshingScopeHandler
}

func GetDefaultStore() store.Store{
	store, err := buntdb.NewStore(":memory:")
	if err != nil {
		logger.Error(err)
		return nil
	}

	return store
}

func (g *Generates) GetAuthorizeToken(ctx context.Context, req *AuthorizeRequest) (*token.Token, error) {
	if fn := g.ClientAuthorizedHandler; fn != nil {
		gt := api.AuthorizationCode
		if req.ResponseType == api.Token {
			gt = api.Implicit
		}

		allowed, err := fn(req.ClientID, gt)
		if err != nil {
			return nil, err
		} else if !allowed {
			return nil, errors.ErrUnauthorizedClient
		}
	}

	if fn := g.ClientScopeHandler; fn != nil {
		allowed, err := fn(req.ClientID, req.Scope)
		if err != nil {
			return nil, err
		} else if !allowed {
			return nil, errors.ErrInvalidScope
		}
	}


	return g.generateAuthToken(ctx, req)
}

func (g *Generates) GetAccessToken(ctx context.Context, tr *TokenRequest) (*token.Token, error)  {
	if allowed := g.checkGrantType(tr.GrantType); !allowed {
		return  nil, errors.ErrUnauthorizedClient
	}

	if fn := g.ClientAuthorizedHandler; fn != nil {
		allowed, err := fn(tr.ClientID, tr.GrantType)
		if err != nil {
			return nil, err
		} else if !allowed {
			return nil, errors.ErrUnauthorizedClient
		}
	}

	switch tr.GrantType {
	case api.AuthorizationCode:
		token, err := g.generateAccessToken(ctx, tr)
		if err != nil {
			switch err {
			case errors.ErrInvalidAuthorizeCode:
				return  nil, errors.ErrInvalidGrant
			case errors.ErrInvalidClient:
				return  nil, errors.ErrInvalidClient
			default:
				return  nil, err
			}
		}

		return token, nil
	case api.PasswordCredentials, api.ClientCredentials:
		if fn := g.ClientScopeHandler; fn != nil {
			allowed, err := fn(tr.ClientID, tr.Scope)
			if err != nil {
				return nil, err
			} else if !allowed {
				return nil, errors.ErrInvalidScope
			}
		}
		return g.generateAccessToken(ctx, tr)
	case api.Refreshing:
		if scope, scopeFn := tr.Scope, g.RefreshingScopeHandler; scope != "" && scopeFn != nil {
			rti, err := g.loadRefreshToken(ctx, tr.Refresh)
			if err != nil {
				if err == errors.ErrInvalidRefreshToken || err == errors.ErrExpiredRefreshToken {
					return nil, errors.ErrInvalidGrant
				}
				return nil, err
			}

			allowed, err := scopeFn(scope, rti.Scope)
			if err != nil {
				return nil, err
			} else if !allowed {
				return nil, errors.ErrInvalidScope
			}
		}

		ti, err := g.refreshAccessToken(ctx, tr)
		if err != nil {
			if err == errors.ErrInvalidRefreshToken || err == errors.ErrExpiredRefreshToken {
				return nil, errors.ErrInvalidGrant
			}
			return nil, err
		}
		return ti, nil
	}

	return nil, errors.ErrUnsupportedGrantType
}

func (g *Generates) LoadAccessToken(ctx context.Context, access string) (*token.Token, error) {
	if access == "" {
		return nil, errors.ErrInvalidAccessToken
	}

	ct := time.Now()
	tk, err := g.TokenStore.GetByAccess(ctx, access)
	if err != nil {
		return nil, err
	} else if tk == nil || tk.Access != access {
		return nil, errors.ErrInvalidAccessToken
	} else if tk.Refresh != "" && tk.RefreshExpiresIn != 0 &&
		tk.RefreshCreateAt.Add(tk.RefreshExpiresIn).Before(ct) {
		return nil, errors.ErrExpiredRefreshToken
	} else if tk.AccessExpiresIn != 0 &&
		tk.AccessCreateAt.Add(tk.AccessExpiresIn).Before(ct) {
		return nil, errors.ErrExpiredAccessToken
	}
	return tk, nil
}

func (g *Generates) generateAuthToken(ctx context.Context, req *AuthorizeRequest) (*token.Token, error) {
	cli, err := g.ClientStore.GetClient(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}else if req.RedirectURI != "" {
		if err := g.ValidateURIHandler(cli.Domain, req.RedirectURI); err != nil {
			return nil, err
		}
	}

	createAt := time.Now()
	token := &token.Token{
		ResponseType:     req.ResponseType,
		ClientID:         req.ClientID,
		UserID:           req.UserID,
		RedirectURI:      req.RedirectURI,
		Scope:            req.Scope,
	}

	switch req.ResponseType {
	case api.Code:
		codeExp := g.codeExp
		if codeExp == 0 {
			codeExp = DefaultCodeExp
		}
		token.CodeCreateAt = createAt
		token.CodeExpiresIn = codeExp

		code, err := g.AuthorizeGenerate.Token(ctx, token)
		if err != nil {
			return nil, err
		}
		token.Code = code
	case api.Token:
		tc := g.tokenConfigs[api.Implicit]
		token.AccessCreateAt = createAt
		token.AccessExpiresIn = tc.AccessTokenExp

		if tc.IsGenerateRefresh {
			token.RefreshCreateAt = createAt
			token.RefreshExpiresIn = tc.RefreshTokenExp
		}

		as, rh, err := g.AccessGenerate.Token(ctx, token, tc.IsGenerateRefresh)
		if err != nil {
			return nil, err
		}
		token.Access = as

		if rh != "" {
			token.Refresh = rh
		}
	}

	if fn := g.ExtensionFieldsHandler; fn != nil {
		token.FieldsValue = fn(token)
	}

	err = g.TokenStore.Create(ctx, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (g *Generates) generateAccessToken(ctx context.Context, tr *TokenRequest) (*token.Token, error) {
	cli, err := g.ClientStore.GetClient(ctx, tr.ClientID)
	if err != nil {
		return nil, err
	}
	if len(cli.Secret) > 0 && tr.ClientSecret != cli.Secret {
		return nil, errors.ErrInvalidClient
	}

	if tr.RedirectURI != "" {
		if err := g.ValidateURIHandler(cli.Domain, tr.RedirectURI); err != nil {
			return nil, err
		}
	}

	if tr.GrantType == api.AuthorizationCode {
		token, err := g.getAndDelAuthorizationCode(ctx, tr)
		if err != nil {
			return nil, err
		}
		tr.UserID = token.UserID
		tr.Scope = token.Scope
	}

	if strings.EqualFold(g.TokenType, "mac") && g.MacAlgorithm == "" {
		return nil, errors.ErrInvalidMacAlgorithm
	}

	tc := g.tokenConfigs[tr.GrantType]
	atx := tc.AccessTokenExp
	createAt := time.Now()
	token := &token.Token{
		TokenType:        g.TokenType,
		GrantType:        tr.GrantType,
		MacAlgorithm: 	  g.MacAlgorithm,
		ClientID:         tr.ClientID,
		UserID:           tr.UserID,
		RedirectURI:      tr.RedirectURI,
		Scope:            tr.Scope,
		AccessCreateAt:   createAt,
		AccessExpiresIn:  atx,
	}

	av, rv, err := g.AccessGenerate.Token(ctx, token, tc.IsGenerateRefresh)
	if err != nil {
		return nil, err
	}
	token.Access = av

	if tc.IsGenerateRefresh {
		token.RefreshCreateAt = createAt
		token.RefreshExpiresIn = tc.RefreshTokenExp
	}
	if rv != "" {
		token.Refresh = rv
	}

	if fn := g.ExtensionFieldsHandler; fn != nil {
		token.FieldsValue = fn(token)
	}

	err = g.TokenStore.Create(ctx, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (g *Generates) getAuthorizationCode(ctx context.Context, code string) (*token.Token, error) {
	token, err := g.TokenStore.GetTokenByCode(ctx, code)
	if err != nil {
		return nil, err
	} else if token == nil || token.Code != code || token.CodeCreateAt.Add(token.CodeExpiresIn).Before(time.Now()) {
		err = errors.ErrInvalidAuthorizeCode
		return nil, errors.ErrInvalidAuthorizeCode
	}
	return token, nil
}

func (g *Generates) getAndDelAuthorizationCode(ctx context.Context, tr *TokenRequest) (*token.Token, error) {
	code := tr.Code
	token, err := g.getAuthorizationCode(ctx, code)
	if err != nil {
		return nil, err
	} else if token.ClientID != tr.ClientID {
		return nil, errors.ErrInvalidAuthorizeCode
	} else if codeURI := token.RedirectURI; codeURI != "" && codeURI != tr.RedirectURI {
		return nil, errors.ErrInvalidAuthorizeCode
	}

	err =  g.TokenStore.RemoveByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (g *Generates) checkGrantType(gt api.GrantType) bool {
	for _, agt := range g.AllowedGrantTypes {
		if agt == gt {
			return true
		}
	}
	return false
}

func (g *Generates) refreshAccessToken(ctx context.Context, tr *TokenRequest) (*token.Token, error) {
	cli, err := g.ClientStore.GetClient(ctx, tr.ClientID)
	if err != nil {
		return nil, err
	} else if tr.ClientSecret != cli.Secret {
		return nil, errors.ErrInvalidClient
	}

	tk, err := g.loadRefreshToken(ctx, tr.Refresh)
	if err != nil {
		return nil, err
	} else if tk.ClientID != tr.ClientID {
		return nil, errors.ErrInvalidRefreshToken
	}

	createAt := time.Now()
	oldAccess, oldRefresh := tk.Access, tk.Refresh

	tk.GrantType = tr.GrantType
	tk.RedirectURI = tr.RedirectURI
	tk.Scope = tr.Scope
	tk.AccessCreateAt = createAt

	rc:= DefaultRefreshTokenCfg
	if v := g.refreshConfig; v != nil {
		rc = v
	}

	if v := rc.AccessTokenExp; v > 0 {
		tk.AccessExpiresIn = v
	}

	if rc.IsResetRefreshTime {
		tk.RefreshCreateAt = createAt
	}

	if v := rc.RefreshTokenExp; v > 0 {
		tk.RefreshExpiresIn = v
	}

	if scope := tr.Scope; scope != "" {
		tk.Scope = scope
	}

	tv, rv, err := g.AccessGenerate.Token(ctx, tk, rc.IsGenerateRefresh)
	if err != nil {
		return nil, err
	}

	tk.Access = tv
	if rv != "" {
		tk.Refresh = rv
	}

	if fn := g.ExtensionFieldsHandler; fn != nil {
		tk.FieldsValue = fn(tk)
	}

	if err := g.TokenStore.Create(ctx, tk); err != nil {
		return nil, err
	}

	if rc.IsRemoveAccess {
		if err := g.TokenStore.RemoveByAccess(ctx, oldAccess); err != nil {
			return nil, err
		}
	}

	if rc.IsRemoveRefreshing && rv != "" {
		if err := g.TokenStore.RemoveByRefresh(ctx, oldRefresh); err != nil {
			return nil, err
		}
	}

	if rv == "" {
		tk.Refresh = ""
		tk.RefreshCreateAt = time.Now()
		tk.RefreshExpiresIn = 0
	}

	return tk, nil
}

func (g *Generates) loadRefreshToken(ctx context.Context, refresh string) (*token.Token, error) {
	if refresh == "" {
		return nil, errors.ErrInvalidRefreshToken
	}

	token, err := g.TokenStore.GetByRefresh(ctx, refresh)
	if err != nil {
		return nil, err
	} else if token == nil || token.Refresh != refresh {
		return nil, errors.ErrInvalidRefreshToken
	} else if token.RefreshExpiresIn != 0 && // refresh token set to not expire
		token.RefreshCreateAt.Add(token.RefreshExpiresIn).Before(time.Now()) {
		return nil, errors.ErrExpiredRefreshToken
	}
	return token, nil
}