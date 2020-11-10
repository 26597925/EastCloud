package server

import (
	"github.com/26597925/EastCloud/pkg/oauth2/api"
	"github.com/26597925/EastCloud/pkg/oauth2/server/errors"
	"github.com/26597925/EastCloud/pkg/oauth2/server/token"
	"net/http"
	"net/url"
	"strings"
)

type (
	ClientInfoHandler            func(r *http.Request) (clientID, clientSecret string, err error)
	PasswordAuthorizationHandler func(username, password string) (userID string, err error)
	UserAuthorizationHandler     func(w http.ResponseWriter, r *http.Request) (userID string, err error)
	AuthorizeScopeHandler        func(w http.ResponseWriter, r *http.Request) (scope string, err error)
	ValidateURIHandler			 func(baseURI, redirectURI string) error
	ExtensionFieldsHandler  	 func(t *token.Token) (fieldsValue map[string]interface{})

	ClientAuthorizedHandler      func(clientID string, grant api.GrantType) (allowed bool, err error)
	ClientScopeHandler			 func(clientID, scope string) (allowed bool, err error)
	RefreshingScopeHandler 		 func(newScope, oldScope string) (allowed bool, err error)

	ValidateClientHandler		 func(nonce string, ts int, ext string) (allowed bool, err error)
)

func ClientFormHandler(r *http.Request) (string, string, error) {
	clientID := r.Form.Get("client_id")
	if clientID == "" {
		return "", "", errors.ErrInvalidClient
	}
	clientSecret := r.Form.Get("client_secret")
	return clientID, clientSecret, nil
}

func ClientBasicHandler(r *http.Request) (string, string, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return "", "", errors.ErrInvalidClient
	}
	return username, password, nil
}

func DefValidateURIHandler(baseURI string, redirectURI string) error {
	base, err := url.Parse(baseURI)
	if err != nil {
		return err
	}

	redirect, err := url.Parse(redirectURI)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(redirect.Host, base.Host) {
		return errors.ErrInvalidRedirectURI
	}
	return nil
}

func DefUserAuthorizationHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	return "", errors.ErrAccessDenied
}

func DefPasswordAuthorizationHandler(username, password string) (string, error) {
	return "", errors.ErrAccessDenied
}