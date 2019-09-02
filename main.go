package main

import (
	"github.com/labstack/echo"
	"github.com/paulohbmatias/oauth2-go/controllers"
	"github.com/subosito/gotenv"
	"log"
	"os"
)

func init() {
	err := gotenv.Load()
	if err != nil{
		log.Fatal(err)
	}
	authController = controllers.NewAuthController()
	port = os.Getenv("PORT")
	e = echo.New()
}

var (
	e *echo.Echo
	authController *controllers.AuthController
	port string
)

func main() {

	g := e.Group("/auth")

	g.POST("/token", authController.TokenController)
	g.POST("/login", authController.Login)
	g.POST("/refreshToken", authController.RefreshToken)
	g.POST("/signUp", authController.SignUp)

	e.Logger.Fatal(e.Start(":"+port))
}
