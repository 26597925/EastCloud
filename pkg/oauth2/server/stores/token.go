package stores

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sapi/pkg/oauth2/server/token"
	"sapi/pkg/store"
	"time"
)

var TokenPrefix = "/sapi/oauth2/token"

type TokenStore struct {
	Store store.Store

	Prefix string
}

func (ts *TokenStore) Create(ctx context.Context, token *token.Token) error {
	ct := time.Now()
	b, err := json.Marshal(token)
	if err != nil {
		return err
	}

	prefix := ts.getPrefix()

	if code := token.Code; code != "" {
		key := fmt.Sprintf("%s/code/%s", prefix, code)
		err := ts.Store.Set(ctx, key, string(b), token.CodeExpiresIn)
		return err
	}

	basicID := uuid.Must(uuid.NewRandom()).String()
	ae := token.AccessExpiresIn
	rp := ae
	if refresh := token.Refresh; refresh != "" {
		rp = token.RefreshCreateAt.Add(token.RefreshExpiresIn).Sub(ct)
		if ae.Seconds() > rp.Seconds() {
			ae = rp
		}
		key := fmt.Sprintf("%s/refresh/%s", prefix, refresh)
		ts.Store.Set(ctx, key, basicID, rp)
		if err != nil {
			return err
		}
	}

	key := fmt.Sprintf("%s/basicID/%s", prefix, basicID)
	err = ts.Store.Set(ctx, key, string(b), rp)
	if err != nil {
		return err
	}

	key = fmt.Sprintf("%s/access/%s", prefix, token.Access)
	err = ts.Store.Set(ctx, key, basicID, ae)
	return err
}

func (ts *TokenStore) GetTokenByCode(ctx context.Context, code string) (*token.Token, error) {
	prefix := ts.getPrefix()
	key := fmt.Sprintf("%s/code/%s", prefix, code)

	return ts.getToken(ctx, key)
}

func (ts *TokenStore) RemoveByCode(ctx context.Context, code string) error {
	prefix := ts.getPrefix()
	key := fmt.Sprintf("%s/code/%s", prefix, code)

	return ts.Store.Delete(ctx, key)
}

func (ts *TokenStore) GetByAccess(ctx context.Context, access string) (*token.Token, error) {
	prefix := ts.getPrefix()
	key := fmt.Sprintf("%s/access/%s", prefix, access)
	basicID, err := ts.Store.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	key = fmt.Sprintf("%s/basicID/%s", prefix, basicID)
	return ts.getToken(ctx, key)
}


func (ts *TokenStore) GetByRefresh(ctx context.Context, refresh string) (*token.Token, error) {
	prefix := ts.getPrefix()

	key := fmt.Sprintf("%s/refresh/%s", prefix, refresh)
	basicID, err := ts.Store.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	key = fmt.Sprintf("%s/basicID/%s", prefix, basicID)
	return ts.getToken(ctx, key)
}

func (ts *TokenStore) RemoveByAccess(ctx context.Context, access string) error {
	prefix := ts.getPrefix()
	key := fmt.Sprintf("%s/access/%s", prefix, access)
	return ts.Store.Delete(ctx, key)
}

func (ts *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	prefix := ts.getPrefix()
	key := fmt.Sprintf("%s/refresh/%s", prefix, refresh)
	return ts.Store.Delete(ctx, key)
}

func (ts *TokenStore) getPrefix() string {
	prefix := ts.Prefix
	if prefix == "" {
		prefix = ClientPrefix
	}

	return prefix
}

func (ts *TokenStore) getToken(ctx context.Context, key string) (*token.Token, error) {
	d, err := ts.Store.Get(ctx, key)
	if err != nil {
		return nil, errors.New("not found")
	}

	var t token.Token
	err = json.Unmarshal([]byte(d), &t)
	if err != nil {
		return nil, errors.New("data parse errors")
	}
	return &t, nil
}