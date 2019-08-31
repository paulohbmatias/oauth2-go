package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/paulohbmatias/oauth2-go/config"
	"github.com/paulohbmatias/oauth2-go/models"
	"github.com/paulohbmatias/oauth2-go/repositories/user"
	"github.com/paulohbmatias/oauth2-go/utils"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/oauth2.v3/errors"
	"log"
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

func (a AuthController) PasswordCredentials(authConfig config.AuthConfig, db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		clientId, secret, _ := r.BasicAuth()

		client, err := authConfig.Manager.GetClient(clientId)

		if err != nil || client.GetSecret() != secret{
			http.Error(w, errors.ErrInvalidClient.Error(), http.StatusInternalServerError)
			return
		}

		var userModel models.User
		var errorModel models.Error
		err = json.NewDecoder(r.Body).Decode(&userModel)
		if err != nil{
			fmt.Println(err)
			return
		}

		if userModel.Email == ""{
			errorModel.Message = "Email is missing"
			utils.SendError(w, http.StatusBadRequest, errorModel)
			return
		}

		if userModel.Password == ""{
			errorModel.Message = "Password is missing"
			utils.SendError(w, http.StatusBadRequest, errorModel)
			return
		}

		password := userModel.Password

		userRepo := user.UserRepository{}
		userModel, err = userRepo.Login(db, userModel)

		if err != nil{
			if err == sql.ErrNoRows{
				errorModel.Message = "The userModel does not exist"
				utils.SendError(w, http.StatusBadRequest, errorModel)
				return
			}else{
				log.Fatal(err)
			}
		}

		hashPassword := userModel.Password

		err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

		if err != nil{
			errorModel.Message = "Invalid password"
			utils.SendError(w, http.StatusBadRequest, errorModel)
			return
		}

		token, err := authConfig.Config.PasswordCredentialsToken(context.TODO(), userModel.Email, userModel.Password)
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




