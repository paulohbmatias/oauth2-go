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
	"github.com/paulohbmatias/oauth2-go/utils"
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
	var errorModel models.Error
	err = json.NewDecoder(c.Request().Body).Decode(&userModel)
	if err != nil {
		errorModel.Message = "Invalid request"
		utils.SendError(c.Response(), http.StatusBadRequest, errorModel)
		return err
	}

	if userModel.Email == "" {
		errorModel.Message = "Email is missing"
		utils.SendError(c.Response(), http.StatusBadRequest, errorModel)
		return err
	}

	if userModel.Password == "" {
		errorModel.Message = "Password is missing"
		utils.SendError(c.Response(), http.StatusBadRequest, errorModel)
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userModel.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	userModel.Password = string(hash)

	userRepo := user.UserRepository{}
	userModel, err = userRepo.SignUp(authController.db, userModel)

	if err != nil {
		fmt.Println(err)
		errorModel.Message = "Server error."
		utils.SendError(c.Response(), http.StatusInternalServerError, errorModel)
		return err
	}

	userModel.Password = ""

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	utils.SendSuccess(c.Response(), userModel)
	return nil
}

func (authController *AuthController) Login(c echo.Context) error{
	clientId, secret, _ := c.Request().BasicAuth()

	client, err := authController.authConfig.Manager.GetClient(clientId)

	if err != nil || client.GetSecret() != secret {
		http.Error(c.Response(), errors.ErrInvalidClient.Error(), http.StatusInternalServerError)
		return err
	}

	var userModel models.User
	var errorModel models.Error


	userModel.Email = c.FormValue("username")
	userModel.Password = c.FormValue("password")

	if userModel.Email == "" {
		errorModel.Message = "Email is missing"
		utils.SendError(c.Response(), http.StatusBadRequest, errorModel)
		return err
	}

	if userModel.Password == "" {
		errorModel.Message = "Password is missing"
		utils.SendError(c.Response(), http.StatusBadRequest, errorModel)
		return err
	}

	password := userModel.Password

	userRepo := user.UserRepository{}
	userModel, err = userRepo.Login(authController.db, userModel)

	fmt.Println(userModel)
	if err != nil {
		if err == sql.ErrNoRows {
			errorModel.Message = "The user does not exist"
			utils.SendError(c.Response(), http.StatusBadRequest, errorModel)
			return nil
		} else {
			log.Fatal(err)
		}
	}

	hashPassword := userModel.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

	if err != nil {
		errorModel.Message = "Invalid password"
		utils.SendError(c.Response(), http.StatusBadRequest, errorModel)
		return nil
	}

	token, err := authController.authConfig.Config.PasswordCredentialsToken(context.Background(), "test", "test")
	if err != nil {
		http.Error(c.Response(), err.Error(), http.StatusInternalServerError)
		return nil
	}

	token.SetAuthHeader(c.Request())

	e := json.NewEncoder(c.Response())
	e.SetIndent("", "  ")
	_ = e.Encode(token)
	return nil
}

func (authController *AuthController) RefreshToken(c echo.Context) error{
	var errorModel models.Error
	authHeader := c.Request().Header.Get("Authorization")
	//bearerToken := strings.Split(authHeader, " ")
	token, err := authController.authConfig.Config.Exchange(context.TODO(), authHeader)

	if err != nil {
		errorModel.Message = "Invalid token"
		utils.SendError(c.Response(), http.StatusBadRequest, errorModel)
	}
	return c.JSON(http.StatusOK, token)
}
