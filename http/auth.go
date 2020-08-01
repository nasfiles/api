package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/color"

	"github.com/nasfiles/nasfilesapi"
)

// AuthLogin authenticates a user
func AuthLogin(w http.ResponseWriter, r *http.Request, c *nasfilesapi.Config) (int, error) {
	credentials := &nasfilesapi.AuthGuess{}

	// Read bytes from request
	reqBuffer := new(bytes.Buffer)
	reqBuffer.ReadFrom(r.Body)

	err := json.Unmarshal(reqBuffer.Bytes(), credentials)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Check credentials size
	if len(credentials.Username) < 6 || len(credentials.Password) < 4 {
		return http.StatusBadRequest, nil
	}

	// Authenticate user
	err = c.Services.Auth.Login(credentials)
	if err != nil {
		fmt.Println(err)
		return http.StatusUnauthorized, err
	}

	// Fetch user information
	user, err := c.Services.User.GetByUsername(credentials.Username)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Create JWT token
	claims := nasfilesapi.APIClaims{
		UID:      user.UID,
		Username: user.Username,
		Email:    user.Email,
		Admin:    user.Admin,
		StandardClaims: jwt.StandardClaims{
			Audience:  "user",
			ExpiresAt: time.Now().Add(time.Hour * 24 * 14).Unix(),
			Issuer:    "nasfilesapi",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(c.PrivateKey))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	res := &nasfilesapi.APITokenResponse{
		Token:      ss,
		Expiration: time.Now().Add(time.Hour * 24 * 14),
	}

	color.HiGreen("User %s just logged in.", credentials.Username)

	return jsonPrint(w, res)
}
