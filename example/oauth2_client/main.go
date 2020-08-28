package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sapi/pkg/logger"
	"sapi/pkg/oauth2/client"
	"time"
)

const (
	authServerURL = "http://localhost:9096"
)

func main() {
	oauth2 := client.NewOauth2(&client.Options{
		ClientID:     "222222",
		ClientSecret: "22222222",
		Endpoint:     client.Endpoint{
			AuthURL:  authServerURL + "/authorize",
			TokenURL: authServerURL + "/token",
		},
		RedirectURL:  "http://localhost:9094/oauth2_client",
		Scopes:       []string{"all"},
		AuthStyle:   client.AuthStyleInHeader,
	})

	var globalToken *client.Token

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u := oauth2.AuthCodeURL("xyz")
		logger.Info(u)
		http.Redirect(w, r, u, http.StatusFound)
	})

	http.HandleFunc("/oauth2_client", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		state := r.Form.Get("state")
		logger.Info(state)
		if state != "xyz" {
			http.Error(w, "State invalid", http.StatusBadRequest)
			return
		}

		error := r.Form.Get("errors")
		if error != "" {
			description := r.Form.Get("error_description")
			info := fmt.Sprintf("%s:%s", error, description)
			http.Error(w, info, http.StatusBadRequest)
			return
		}

		code := r.Form.Get("code")
		logger.Info(code)
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}

		token, err := oauth2.AuthorizationCodeToken(context.Background(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		globalToken = token
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)

	})

	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		if globalToken == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		globalToken.Expiry = time.Now()
		token, err := oauth2.TokenSource(context.Background(), globalToken).Token()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		globalToken = token
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	http.HandleFunc("/try", func(w http.ResponseWriter, r *http.Request) {
		if globalToken == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		resp, err := http.Get(fmt.Sprintf("%s/test?access_token=%s", authServerURL, globalToken.AccessToken))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		io.Copy(w, resp.Body)
	})

	http.HandleFunc("/try1", func(w http.ResponseWriter, r *http.Request) {
		if globalToken == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		globalToken.SetTs(time.Now().Unix())
		globalToken.SetNonce("123456")
		urlStr := fmt.Sprintf("%s/test1", authServerURL)
		resp, err := oauth2.Client(context.Background(), globalToken).Get(urlStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		io.Copy(w, resp.Body)
	})

	http.HandleFunc("/pwd", func(w http.ResponseWriter, r *http.Request) {
		token, err := oauth2.PasswordCredentialsToken(context.Background(), "test", "test")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		globalToken = token
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	http.HandleFunc("/client", func(w http.ResponseWriter, r *http.Request) {
		token, err := oauth2.ClientCredentialsToken(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	logger.Info("Client is running at 9094 port.Please open http://localhost:9094")
	http.ListenAndServe(":9094", nil)
}
