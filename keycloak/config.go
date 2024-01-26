package keycloak

type JWT struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
	IDToken          string `json:"id_token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type KeycloakConfig struct {
	Address      string
	Port         string
	ClientId     string // clientId specified in Keycloak
	ClientSecret string // client secret specified in Keycloak
	Realm        string // realm specified in Keycloak
	Protocol     string
}

func GetNewKeycloakConfig() *KeycloakConfig {
	return &KeycloakConfig{
		Address:      "192.168.86.211",
		Port:         "32088",
		ClientId:     "auth-svc",
		ClientSecret: "3nUB5EUG3fenYFa9xqnF376PLFZHWxFV",
		Realm:        "nshub",
		Protocol:     "openid-connect",
	}
}
