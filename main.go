package main

import (
	"database/sql"
	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
	"github.com/paulohbmatias/oauth2-go/config"
	"github.com/paulohbmatias/oauth2-go/controllers"
	"github.com/paulohbmatias/oauth2-go/driver"
	"github.com/paulohbmatias/oauth2-go/middlewares"
	"github.com/subosito/gotenv"
	"log"
	"os"
)

func init() {
	err := gotenv.Load()
	if err != nil{
		log.Fatal(err)
	}
	db = driver.ConnectDB()
	authConfig = config.NewAuthConfig(db)
	authController = controllers.NewAuthController(authConfig, db)
	middleware = middlewares.NewMiddleware(authConfig)
	port = os.Getenv("PORT")
	e = echo.New()
}

var (
	e *echo.Echo
	db *sql.DB
	authConfig *config.AuthConfig
	authController *controllers.AuthController
	middleware *middlewares.Middleware
	port string
)

func main() {

	g := e.Group("/auth")
	g.Use(echoMiddleware.BasicAuth(middleware.BasicAuth))

	g.POST("/token", authController.TokenController)
	g.POST("/login", authController.Login)
	g.POST("/refreshToken", authController.RefreshToken)
	g.POST("/signUp", authController.SignUp)

	e.Logger.Fatal(e.Start(":"+port))
}
