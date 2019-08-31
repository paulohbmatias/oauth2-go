package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/paulohbmatias/oauth2-go/config"
	"github.com/paulohbmatias/oauth2-go/controllers"
	"github.com/paulohbmatias/oauth2-go/driver"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
)

func init() {
	err := gotenv.Load()
	if err != nil{
		log.Fatal(err)
	}
	authConfig.SetupConfig()
	authConfig.SetupClients()
	authConfig.SetupManager()
	authConfig.SetupServer()
}

var (
	authConfig config.AuthConfig
	authController controllers.AuthController
	db *sql.DB
)

func main() {
	db = driver.ConnectDB()
	router := mux.NewRouter()

	router.HandleFunc("/oauth2/token", authController.TokenController(authConfig))
	router.HandleFunc("/oauth2/resourceOwner", authController.PasswordCredentials(authConfig, db)).Methods("POST")
	router.HandleFunc("/oauth2/refreshToken", authController.PasswordCredentials(authConfig)).Methods("POST")

	log.Println("Server is running at 9096 port.")
	log.Fatal(http.ListenAndServe(":9096", router))
}
