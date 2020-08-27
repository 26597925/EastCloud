package generates

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/google/uuid"
	"sapi/pkg/oauth2/server/token"
	"strconv"
	"strings"
)

func NewBaseAccess() *BaseAccess {
	return &BaseAccess{}
}

type BaseAccess struct {
}

// Token based on the UUID generated token
func (ag *BaseAccess) Token(ctx context.Context, t *token.Token, isGenRefresh bool) (string, string, error) {
	buf := bytes.NewBufferString(t.ClientID)
	buf.WriteString(t.UserID)
	buf.WriteString(strconv.FormatInt(t.AccessCreateAt.UnixNano(), 10))

	access := base64.URLEncoding.EncodeToString([]byte(uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes()).String()))
	access = strings.ToUpper(strings.TrimRight(access, "="))
	refresh := ""
	if isGenRefresh {
		refresh = base64.URLEncoding.EncodeToString([]byte(uuid.NewSHA1(uuid.Must(uuid.NewRandom()), buf.Bytes()).String()))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	if strings.EqualFold(t.TokenType, "mac") {
		key := base64.URLEncoding.EncodeToString([]byte(uuid.NewSHA1(uuid.Must(uuid.NewRandom()), buf.Bytes()).String()))
		key = strings.ToUpper(strings.TrimRight(refresh, "="))
		t.MacKey = key
	}

	return access, refresh, nil
}

