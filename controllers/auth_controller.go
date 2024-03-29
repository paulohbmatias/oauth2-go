package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulohbmatias/oauth2-go/config"
	"github.com/paulohbmatias/oauth2-go/models"
	"github.com/paulohbmatias/oauth2-go/repositories/user"
	"golang.org/x/crypto/bcrypt"
	oauth22 "golang.org/x/oauth2"
	"gopkg.in/oauth2.v3"
	"log"
	"net/http"
	"time"
)

type AuthController struct {
	AuthConfig *config.AuthConfig
	db         *sql.DB
}

func NewAuthController(authConfig *config.AuthConfig, db *sql.DB) *AuthController {
	return &AuthController{AuthConfig: authConfig, db: db}
}


func (authController *AuthController) TokenController(c echo.Context) error {
	err := authController.AuthConfig.Server.HandleTokenRequest(c.Response(), c.Request())
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return nil
	}
	return nil
}

func (authController *AuthController) SignUp(c echo.Context) error{

	var userModel models.User
	err := json.NewDecoder(c.Request().Body).Decode(&userModel)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	if userModel.Email == "" {
		return c.String(http.StatusBadRequest, "Email is missing")
	}

	if userModel.Password == "" {
		return c.String(http.StatusBadRequest, "Password is missing")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userModel.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	userModel.Password = string(hash)

	userRepo := user.UserRepository{}
	userModel, err = userRepo.SignUp(authController.db, userModel)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Server error.")
	}
	userModel.Password = ""

	return c.JSON(http.StatusCreated, userModel)
}

func (authController *AuthController) Login(c echo.Context) error{
	var userModel models.User

	userModel.Email = c.FormValue("username")
	userModel.Password = c.FormValue("password")

	token, err := authController.AuthConfig.Config.PasswordCredentialsToken(context.TODO(), userModel.Email, userModel.Password)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, token)
}

func (authController *AuthController) RefreshToken(c echo.Context) error{
	clientId, secret, _ := c.Request().BasicAuth()

	refreshToken := c.FormValue("refresh_token")

	tokenInfo, err := authController.AuthConfig.Manager.LoadRefreshToken(refreshToken)

	if err != nil{
		fmt.Println(err.Error())
		return c.String(http.StatusBadRequest, "Invalid token")
	}

	tokenGenerate := &oauth2.TokenGenerateRequest{
		ClientID:       clientId,
		ClientSecret:   secret,
		UserID:         tokenInfo.GetUserID(),
		RedirectURI:    tokenInfo.GetRedirectURI(),
		Scope:          tokenInfo.GetScope(),
		Code:           tokenInfo.GetCode(),
		Refresh:        tokenInfo.GetRefresh(),
		AccessTokenExp: tokenInfo.GetAccessExpiresIn(),
		Request:        c.Request(),
	}


	newToken, err := authController.AuthConfig.Manager.RefreshAccessToken(tokenGenerate)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid token")
	}

	token := oauth22.Token{
		AccessToken:  newToken.GetAccess(),
		TokenType:    "Bearer",
		RefreshToken: newToken.GetRefresh(),
		Expiry:       time.Now().Add(newToken.GetAccessExpiresIn()),
	}

	return c.JSON(http.StatusOK, token)
}
