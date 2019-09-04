package middlewares

import (
	"github.com/labstack/echo"
	"github.com/paulohbmatias/oauth2-go/config"
)

type Middleware struct {
	authConfig *config.AuthConfig
}

func NewMiddleware(authConfig *config.AuthConfig) *Middleware {
	return &Middleware{authConfig: authConfig}
}

func (m *Middleware) BasicAuth(clientId, secret string, c echo.Context) (bool, error) {
	client, err := m.authConfig.Manager.GetClient(clientId)

	if err != nil || client.GetSecret() != secret {
		return false, nil
	}
	return true, nil
}
