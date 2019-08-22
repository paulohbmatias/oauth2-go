package config

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/server"
	"log"
	"oauth/oauth2/generates"
	"oauth/oauth2/manage"
	"oauth/oauth2/models"
	"oauth/oauth2/store"
	"os"
)

type AuthConfig struct {
	Server  *server.Server
	Config  oauth2.Config
	Manager *manage.Manager
	ClientStore *store.ClientStore
}

const (
	authServerURL = "http://localhost:9096"
)

//func (a *AuthConfig) GetAuthConfig() *AuthConfig{
//	auth := AuthConfig{}
//
//	auth.s
//}

func (a *AuthConfig) SetupManager(){
	a.Manager = manage.NewDefaultManager()
	a.Manager.SetAuthorizeCodeTokenCfg(manage.DefaultPasswordTokenCfg)

	// token store
	a.Manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	a.Manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte(os.Getenv("PRIVATE_SECRET")), jwt.SigningMethodHS512))

	a.Manager.MapClientStorage(a.ClientStore)
}

func (a *AuthConfig) SetupServer() {
	a.Server = server.NewServer(server.NewConfig(), a.Manager)

	a.Server.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		if username == "test" && password == "test" {
			userID = "test"
		}
		return
	})

	a.Server.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	a.Server.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})
}

func (a *AuthConfig) SetupClients(){
	a.ClientStore = store.NewClientStore()
	err := a.ClientStore.Set(os.Getenv("CLIENT_ID"), &models.Client{
		ID:     os.Getenv("CLIENT_ID"),
		Secret: os.Getenv("CLIENT_SECRET"),
	})
	if err != nil{
		fmt.Println(err)
		return
	}
}

func (a *AuthConfig) SetupConfig(){
	a.Config = oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"all"},
		Endpoint: oauth2.Endpoint{
			TokenURL: authServerURL + "/token",
		},
	}
}