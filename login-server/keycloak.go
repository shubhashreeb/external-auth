package loginserver

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"gitlab.com/sitenet/svclib/logger"
)

type KeycloakConfig struct {
	address      string
	clientId     string // clientId specified in Keycloak
	clientSecret string // client secret specified in Keycloak
	realm        string // realm specified in Keycloak
}

type Keycloak struct {
	client *gocloak.GoCloak
	config KeycloakConfig
	log    logger.Logger
}

func NewKeycloak(log logger.Logger) *Keycloak {
	c := KeycloakConfig{
		address:      "http://192.168.86.211:32088/",
		clientId:     "auth-svc",
		clientSecret: "PIvkN94ImvghqZmv0vJUO2PElHtWYXsY", //"eQxQdPudNTyi8rv5L3Tgs1SO5byD6vNB",
		realm:        "nshub",
	}
	return &Keycloak{
		config: c,
		client: gocloak.NewClient(c.address),
		log:    log,
	}
}

func (k *Keycloak) GetLoginToken(req LoginRequest) *LoginResponse {
	jwt, err := k.client.Login(context.Background(),
		k.config.clientId,
		k.config.clientSecret,
		k.config.realm,
		req.Username,
		req.Password,
	)

	fmt.Println("Login request received", jwt, err)

	if err != nil {
		// http.Error(w, err.Error(), http.StatusForbidden)
		return &LoginResponse{}
	}
	fmt.Println("Here is the token response", jwt)

	// TODO - Add to cache

	return &LoginResponse{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
		ExpiresIn:    jwt.ExpiresIn,
	}

}

func (k *Keycloak) Logout(req LogoutRequest) string {
	fmt.Println("Logout req received with refresh token", req.RefreshToken)

	err := k.client.Logout(context.Background(),
		k.config.clientId,
		k.config.clientSecret,
		k.config.realm,
		req.RefreshToken,
	)
	fmt.Println("Logout response received", err)

	if err != nil {
		// http.Error(w, err.Error(), http.StatusForbidden)
		return fmt.Sprintf("invalid with error", err)
	}
	// fmt.Println("Here is the token response", jwt)

	// TODO - Add to cache

	return "Sucessfully logged out"
}
