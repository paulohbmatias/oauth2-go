package main

import (
	"github.com/gorilla/mux"
	"github.com/paulohbmatias/oauth2/config"
	"github.com/paulohbmatias/oauth2/controllers"
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
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/token", controllers.AuthController{}.TokenController(authConfig))

	router.HandleFunc("/pwd", controllers.AuthController{}.PasswordCredentials(authConfig)).Methods("POST")

	log.Println("Server is running at 9096 port.")
	log.Fatal(http.ListenAndServe(":9096", router))
}
