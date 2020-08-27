package server

import (
	"net/http"
	"sapi/pkg/oauth2/api"
	"sapi/pkg/oauth2/server/errors"
	"sapi/pkg/util/stringext"
	"strconv"
	"strings"
)

type AuthorizeRequest struct {
	ResponseType   api.ResponseType
	ClientID       string
	ClientSecret   string
	Scope          string
	RedirectURI    string
	State          string

	UserID         string
}

type TokenRequest struct {
	GrantType	   api.GrantType
	ClientID       string
	ClientSecret   string
	Code           string
	RedirectURI    string

	Scope          string
	Refresh        string
	UserID         string
}

type MacRequest struct {
	id 		string
	ts 		int
	nonce 	string //客户端生成的唯一值
	ext 	string
	mac 	string
}

type Request struct {
	Request *http.Request
	AllowedResponseTypes  []api.ResponseType
	AllowGetAccessRequest bool

	ClientInfoHandler   ClientInfoHandler
	PasswordAuthorizationHandler PasswordAuthorizationHandler
}

func (re *Request) checkResponseType(rt api.ResponseType) bool {
	for _, art := range re.AllowedResponseTypes {
		if art == rt {
			return true
		}
	}
	return false
}

func (re *Request) ValidationAuthorizeRequest() (*AuthorizeRequest, error) {
	redirectURI := re.Request.FormValue("redirect_uri")
	clientID := re.Request.FormValue("client_id")
	if !(re.Request.Method == "GET" || re.Request.Method == "POST") ||
		clientID == "" {
		return nil, errors.ErrInvalidRequest
	}

	resType := api.ResponseType(re.Request.FormValue("response_type"))
	if resType.String() == "" {
		return nil, errors.ErrUnsupportedResponseType
	} else if allowed := re.checkResponseType(resType); !allowed {
		return nil, errors.ErrUnauthorizedClient
	}

	req := &AuthorizeRequest{
		RedirectURI:  redirectURI,
		ResponseType: resType,
		ClientID:     clientID,
		State:        re.Request.FormValue("state"),
		Scope:        re.Request.FormValue("scope"),
	}

	return req, nil
}

func (re *Request) ValidationTokenRequest() (*TokenRequest, error) {
	if v := re.Request.Method; !(v == "POST" ||
		(re.AllowGetAccessRequest && v == "GET")) {
		return nil, errors.ErrInvalidRequest
	}

	gt := api.GrantType(re.Request.FormValue("grant_type"))
	if gt.String() == "" {
		return nil, errors.ErrUnsupportedGrantType
	}

	clientID, clientSecret, err := re.ClientInfoHandler(re.Request)
	if err != nil {
		return nil, err
	}

	tr := &TokenRequest{
		GrantType:    gt,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	switch gt {
	case api.AuthorizationCode:
		tr.RedirectURI = re.Request.FormValue("redirect_uri")
		tr.Code = re.Request.FormValue("code")
		if tr.RedirectURI == "" ||
			tr.Code == "" {
			return nil, errors.ErrInvalidRequest
		}
	case api.PasswordCredentials:
		tr.Scope = re.Request.FormValue("scope")
		username, password := re.Request.FormValue("username"), re.Request.FormValue("password")
		if username == "" || password == "" {
			return nil, errors.ErrInvalidRequest
		}

		userID, err := re.PasswordAuthorizationHandler(username, password)
		if err != nil {
			return nil, err
		} else if userID == "" {
			return nil, errors.ErrInvalidGrant
		}
		tr.UserID = userID
	case api.ClientCredentials:
		tr.Scope = re.Request.FormValue("scope")
	case api.Refreshing:
		tr.Refresh = re.Request.FormValue("refresh_token")
		tr.Scope = re.Request.FormValue("scope")
		if tr.Refresh == "" {
			return nil, errors.ErrInvalidRequest
		}
	}

	return tr, nil
}

func (re *Request) BearerAuth() (string, bool) {
	auth := re.Request.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = re.Request.FormValue("access_token")
	}

	return token, token != ""
}

//id / ts / nonce / ext / mac
func (re *Request) MacAuth() (*MacRequest, error) {
	auth := re.Request.Header.Get("Authorization")
	prefix := "MAC "
	token := ""
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}

	mr := &MacRequest{}
	s := strings.Split(token, ":")
	if len(s) == 5 {
		ts,err := strconv.Atoi(s[1])
		if err != nil {
			return nil, err
		}
		mr.id = s[0]
		mr.ts = ts
		mr.nonce = s[2]
		mr.ext = s[3]
		mr.mac = s[4]
	} else {
		s = strings.Split(token, ", ")
		for _, val := range s {
			if strings.HasPrefix(val, "id=") {
				mr.id = stringext.TrimPreVal(val,"id=")
			}
			if strings.HasPrefix(val, "ts=") {
				ts, err := strconv.Atoi(val[len("ts="):])
				if err != nil {
					return nil, err
				}
				mr.ts = ts
			}
			if strings.HasPrefix(val, "nonce=") {
				mr.nonce = stringext.TrimPreVal(val,"nonce=")
			}
			if strings.HasPrefix(val, "ext=") {
				mr.ext = stringext.TrimPreVal(val,"ext=")
			}
			if strings.HasPrefix(val, "mac=") {
				mr.mac = stringext.TrimPreVal(val,"mac=")
			}
		}
	}

	return mr, nil
}