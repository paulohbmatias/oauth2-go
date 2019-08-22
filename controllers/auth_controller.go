package controllers

import (
	"context"
	"encoding/json"
	"github.com/paulohbmatias/oauth2/config"
	"net/http"
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




