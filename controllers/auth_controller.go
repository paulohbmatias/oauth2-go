package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/paulohbmatias/oauth2-go/config"
	"github.com/paulohbmatias/oauth2-go/driver"
	"github.com/paulohbmatias/oauth2-go/models"
	"github.com/paulohbmatias/oauth2-go/repositories/user"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/oauth2.v3/errors"
	"log"
	"net/http"
)

type AuthController struct {
	authConfig *config.AuthConfig
	db         *sql.DB
}

func NewAuthController() *AuthController {
	authConfig := config.NewAuthConfig()
	db := driver.ConnectDB()
	return &AuthController{
		authConfig: authConfig,
		db:         db,
	}
}

func (authController *AuthController) TokenController(c echo.Context) error {
	err := authController.authConfig.Server.HandleTokenRequest(c.Response(), c.Request())
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return nil
	}
	return nil
}

func (authController *AuthController) SignUp(c echo.Context) error{
	clientId, secret, _ := c.Request().BasicAuth()

	client, err := authController.authConfig.Manager.GetClient(clientId)

	if err != nil || client.GetSecret() != secret {
		http.Error(c.Response(), errors.ErrInvalidClient.Error(), http.StatusInternalServerError)
		return err
	}

	var userModel models.User
	err = json.NewDecoder(c.Request().Body).Decode(&userModel)

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
	clientId, secret, _ := c.Request().BasicAuth()

	client, err := authController.authConfig.Manager.GetClient(clientId)

	if err != nil || client.GetSecret() != secret {
		http.Error(c.Response(), errors.ErrInvalidClient.Error(), http.StatusInternalServerError)
		return err
	}

	var userModel models.User

	userModel.Email = c.FormValue("username")
	userModel.Password = c.FormValue("password")


	if userModel.Email == "" {
		return c.String(http.StatusBadRequest, "Email is missing")
	}

	if userModel.Password == "" {
		return c.String(http.StatusBadRequest, "Password is missing")
	}

	password := userModel.Password

	userRepo := user.UserRepository{}
	userModel, err = userRepo.Login(authController.db, userModel)

	fmt.Println(userModel)
	if err != nil {
		return c.String(http.StatusBadRequest, "The user does not exist")
	}

	hashPassword := userModel.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid password")
	}

	token, err := authController.authConfig.Config.PasswordCredentialsToken(context.Background(), userModel.Email, userModel.Password)
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, token)
}

func (authController *AuthController) RefreshToken(c echo.Context) error{
	authHeader := c.Request().Header.Get("Authorization")
	//bearerToken := strings.Split(authHeader, " ")
	token, err := authController.authConfig.Config.Exchange(context.TODO(), authHeader)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid token")
	}
	return c.JSON(http.StatusOK, token)
}
