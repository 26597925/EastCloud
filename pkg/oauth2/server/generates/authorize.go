package generates

import (
	"bytes"
	"context"
	"encoding/base64"
	"github.com/26597925/EastCloud/pkg/oauth2/server/token"
	"strings"

	"github.com/google/uuid"
)

func NewAuthorize() *Authorize {
	return &Authorize{}
}

type Authorize struct{}

func (ag *Authorize) Token(ctx context.Context, t *token.Token) (string, error) {
	buf := bytes.NewBufferString(t.ClientID)
	buf.WriteString(t.UserID)
	token := uuid.NewMD5(uuid.Must(uuid.NewRandom()), buf.Bytes())
	code := base64.URLEncoding.EncodeToString([]byte(token.String()))
	code = strings.ToUpper(strings.TrimRight(code, "="))

	return code, nil
}
