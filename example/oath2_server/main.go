package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-session/session"
	"log"
	"net/http"
	"net/url"
	"os"
	"sapi/pkg/oauth2/api"
	"sapi/pkg/oauth2/server"
	"sapi/pkg/oauth2/server/generates"
	"sapi/pkg/oauth2/server/stores"
	"time"
)

var oauth *server.Oauth2

func main() {

	oauth = server.NewOauth2(&server.Options{
		TokenType:            "Mac",//Bearer,Mac
		MacAlgorithm:		  "hmac-sha-1",//hmac-sha-256
		AllowedResponseTypes: []api.ResponseType{api.Code, api.Token},
		AllowedGrantTypes: []api.GrantType{
			api.AuthorizationCode,
			api.PasswordCredentials,
			api.ClientCredentials,
			api.Refreshing,
		},
	})

	oauth.SaveClient(context.Background(), &stores.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:9094",
		UserID: "asd",
	})

	oauth.SetAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	oauth.UserAuthorizationHandler = userAuthorizeHandler

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth1", authHandler)

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		store, err := session.Start(r.Context(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var form url.Values
		if v, ok := store.Get("ReturnUri"); ok {
			form = v.(url.Values)
		}
		r.Form = form

		store.Delete("ReturnUri")
		store.Save()

		err = oauth.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		err := oauth.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		token, err := oauth.ValidationBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(token)

		data := map[string]interface{}{
			"expires_in": int64(token.AccessCreateAt.Add(token.AccessExpiresIn).Sub(time.Now()).Seconds()),
			"client_id":  token.ClientID,
			"user_id":    token.UserID,
		}
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(data)
	})

	http.HandleFunc("/test1", func(w http.ResponseWriter, r *http.Request) {
		token, err := oauth.ValidationMACToken(r)
		fmt.Println(err)
		fmt.Println(token)

		data := map[string]interface{}{
			"expires_in": int64(token.AccessCreateAt.Add(token.AccessExpiresIn).Sub(time.Now()).Seconds()),
			"client_id":  token.ClientID,
			"user_id":    token.UserID,
		}
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(data)
	})

	log.Println("Server is running at 9096 port.")
	log.Fatal(http.ListenAndServe(":9096", nil))
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		store.Set("LoggedInUserID", r.Form.Get("username"))
		store.Save()

		w.Header().Set("Location", "/auth1")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "static/login.html")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	outputHTML(w, r, "static/auth1.html")
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}