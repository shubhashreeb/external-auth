package keycloak

type JWT struct {
	AccessToken      string `json:"access_token"`
	IDToken          string `json:"id_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
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
