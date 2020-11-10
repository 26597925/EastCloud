package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/26597925/EastCloud/pkg/util/crypto"
	"golang.org/x/net/context/ctxhttp"
	"io"
	"io/ioutil"
	"math"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const expiryDelta = 10 * time.Second

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Expiry time.Time `json:"expiry,omitempty"`
	MacKey string `json:"mac_key,omitempty"`
	MacAlgorithm string `json:"mac_algorithm,omitempty"`
	raw interface{}
	ts	int64
	nonce string
	ext string
}

func (t *Token) Type() string {
	if strings.EqualFold(t.TokenType, "bearer") {
		return "Bearer"
	}
	if strings.EqualFold(t.TokenType, "mac") {
		return "MAC"
	}
	if strings.EqualFold(t.TokenType, "basic") {
		return "Basic"
	}
	if t.TokenType != "" {
		return t.TokenType
	}
	return "Bearer"
}

func (t *Token) SetAuthHeader(r *http.Request) {
	if t.Type() == "mac" {
		r.Header.Set("Authorization", t.Type()+" "+t.AccessToken)
	} else {
		r.Header.Set("Authorization", t.buildAuthorization(r))
	}
}

func (t *Token) Valid() bool {
	return t != nil && t.AccessToken != "" && !t.expired()
}

func (t *Token) SetTs(ts int64) {
	t.ts = ts
}

func (t *Token) SetNonce(nonce string) {
	t.nonce = nonce
}

func (t *Token) SetExt(ext string) {
	t.ext = ext
}

func (t *Token) buildAuthorization(r *http.Request) string {
	text := fmt.Sprintf("%d\n%s\n%s\n%s\n%s\n%s\n%s\n", t.ts, t.nonce, r.Method, r.URL.RequestURI(), r.URL.Hostname(), r.URL.Port(), t.ext)
	fmt.Println(text)
	macText := ""
	if strings.EqualFold(t.MacAlgorithm, "hmac-sha-1") {
		macText = crypto.HmacSha1(t.MacKey, text)
	}
	if strings.EqualFold(t.MacAlgorithm, "hmac-sha-256") {
		macText = crypto.HmacSha256(t.MacKey, text)
	}
	return fmt.Sprintf("MAC %s:%d:%s:%s:%s", t.AccessToken, t.ts, t.nonce, t.ext, macText)
}

func (t *Token) expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Round(0).Add(-expiryDelta).Before(time.Now())
}

type tokenJSON struct {
	AccessToken  string         `json:"access_token"`
	TokenType    string         `json:"token_type"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    expirationTime `json:"expires_in"` // at least PayPal returns string, while most return number
	MacKey string 				`json:"mac_key"`
	MacAlgorithm string 		`json:"mac_algorithm"`
}

func (e *tokenJSON) expiry() (t time.Time) {
	if v := e.ExpiresIn; v != 0 {
		return time.Now().Add(time.Duration(v) * time.Second)
	}
	return
}

type expirationTime int32

func (e *expirationTime) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	i, err := n.Int64()
	if err != nil {
		return err
	}
	if i > math.MaxInt32 {
		i = math.MaxInt32
	}
	*e = expirationTime(i)
	return nil
}

func cloneURLValues(v url.Values) url.Values {
	v2 := make(url.Values, len(v))
	for k, vv := range v {
		v2[k] = append([]string(nil), vv...)
	}
	return v2
}

func RequestToken(ctx context.Context, opt *Options, v url.Values) (*Token, error) {
	if opt.AuthStyle != AuthStyleInHeader {
		v = cloneURLValues(v)
		if opt.ClientID != "" {
			v.Set("client_id", opt.ClientID)
		}
		if opt.ClientSecret != "" {
			v.Set("client_secret", opt.ClientSecret)
		}
	}

	req, err := http.NewRequest("POST", opt.Endpoint.TokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if opt.AuthStyle== AuthStyleInHeader {
		req.SetBasicAuth(url.QueryEscape(opt.ClientID), url.QueryEscape(opt.ClientSecret))
	}

	r, err := ctxhttp.Do(ctx, ContextClient(ctx), req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	r.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("oauth2_client: cannot fetch token: %v", err)
	}
	if code := r.StatusCode; code < 200 || code > 299 {
		return nil, fmt.Errorf("oauth2_client: cannot fetch token: %v\nResponse: %s", r.Status, r.Body)
	}

	fmt.Println(string(body))

	var token *Token
	content, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch content {
	case "application/x-www-form-urlencoded", "text/plain":
		vals, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		token = &Token{
			AccessToken:  vals.Get("access_token"),
			TokenType:    vals.Get("token_type"),
			RefreshToken: vals.Get("refresh_token"),
			MacKey:		  vals.Get("mac_key"),
			MacAlgorithm: vals.Get("mac_algorithm"),
			raw:          vals,
		}
		e := vals.Get("expires_in")
		expires, _ := strconv.Atoi(e)
		if expires != 0 {
			token.Expiry = time.Now().Add(time.Duration(expires) * time.Second)
		}
	default:
		var tj tokenJSON
		if err = json.Unmarshal(body, &tj); err != nil {
			return nil, err
		}
		token = &Token{
			AccessToken:  tj.AccessToken,
			TokenType:    tj.TokenType,
			RefreshToken: tj.RefreshToken,
			Expiry:       tj.expiry(),
			MacKey:		  tj.MacKey,
			MacAlgorithm: tj.MacAlgorithm,
			raw:          make(map[string]interface{}),
		}
		json.Unmarshal(body, &token.raw) // no errors checks for optional fields
	}
	if token.AccessToken == "" {
		return nil, errors.New("oauth2_client: server response missing access_token")
	}
	return token, nil
}
