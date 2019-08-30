package controllers

import (
	"context"
	"encoding/json"
	"github.com/paulohbmatias/oauth2/config"
	"net/http"
	"oauth/oauth2/errors"
	"os"
)

type AuthController struct {}

func (a AuthController) TokenController(authConfig config.AuthConfig) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){
		err := authConfig.Server.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (a AuthController) PasswordCredentials(authConfig config.AuthConfig) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		email, password, ok := r.BasicAuth()
		if !(email == os.Getenv("CLIENT_ID") && password == os.Getenv("CLIENT_SECRET") && ok){
			http.Error(w, errors.ErrAccessDenied.Error(), http.StatusInternalServerError)
			return
		}

		token, err := authConfig.Config.PasswordCredentialsToken(context.TODO(), "test", "test")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token.SetAuthHeader(r)

		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		_ = e.Encode(token)
	}
}




