package main

import (
	"context"
	"sapi/pkg/logger"
	"sapi/pkg/oauth2/api"
	"sapi/pkg/oauth2/client"
)

const (
	authServerURL = "http://localhost:9096"
)

func main() {
	oauth2 := client.NewOauth2(&client.Options{
		ClientID:     "222222",
		ClientSecret: "22222222",
		Endpoint:     api.Endpoint{
			AuthURL:  authServerURL + "/authorize",
			TokenURL: authServerURL + "/token",
		},
		RedirectURL:  "http://localhost:9094/oauth2_client",
		Scopes:       []string{"all"},
		AuthStyle:   api.AuthStyleInHeader,
	})

	token, err := oauth2.ClientCredentialsToken(context.Background())
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info(token)
	}
}

