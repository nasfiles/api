package nasfilesapi

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthGuess is used to authenticate a user
type AuthGuess struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// APIClaims is a structure containing the information
// which will be stored in the token payload
type APIClaims struct {
	UID      int    `json:"UID"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Admin    bool   `json:"Adm"`
	jwt.StandardClaims
}

// APITokenResponse is a structure used to send the token to the user
// when him/her logs in
type APITokenResponse struct {
	Token      string
	Expiration time.Time
}

// AuthService ...
type AuthService interface {
	Login(guess *AuthGuess) error
}
