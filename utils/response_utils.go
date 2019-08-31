package utils


import (
	"encoding/json"
	"fmt"
	"github.com/paulohbmatias/oauth2-go/models"
	"net/http"
)

func SendError(w http.ResponseWriter, status int, error models.Error) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(error)
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
	fmt.Println(data)
	_ = json.NewEncoder(w).Encode(data)
}