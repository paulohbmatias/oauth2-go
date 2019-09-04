package config

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/paulohbmatias/oauth2-go/models"
	userRepository "github.com/paulohbmatias/oauth2-go/repositories/user"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"log"
	authModels "oauth/oauth2/models"
	"os"
)

type AuthConfig struct {
	Server  *server.Server
	Config  *oauth2.Config
	Manager *manage.Manager
	ClientStore *store.ClientStore
	db *sql.DB
}

func NewAuthConfig(db *sql.DB) *AuthConfig {
	authConfig := AuthConfig{
		db: db,
	}
	authConfig.SetupConfig()
	authConfig.SetupClients()
	authConfig.SetupManager()
	authConfig.SetupServer()
	return &authConfig
}

const (
	authServerURL = "http://localhost:9096"
)

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

		if username == "" || password == ""{
			err = errors.ErrAccessDenied
			return
		}

		var user models.User

		user.Password = password
		user.Email = username

		userRepo := userRepository.UserRepository{}
		userModel, err := userRepo.Login(a.db, user)

		if err != nil {
			return
		}

		hashPassword := userModel.Password

		err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

		if err != nil{
			return
		}

		userID = userModel.ID

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
	err := a.ClientStore.Set(os.Getenv("CLIENT_ID"), &authModels.Client{
		ID:     os.Getenv("CLIENT_ID"),
		Secret: os.Getenv("CLIENT_SECRET"),
	})
	if err != nil{
		fmt.Println(err)
		return
	}
}

func (a *AuthConfig) SetupConfig(){
	a.Config = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"all"},
		Endpoint: oauth2.Endpoint{
			TokenURL: authServerURL + "/auth/token",
		},
	}
}