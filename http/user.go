package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"

	"github.com/nasfiles/nasfilesapi"
	"github.com/nasfiles/nasfilesapi/utils"
)

//UserAdd inserts a user into the database
func UserAdd(w http.ResponseWriter, r *http.Request, c *nasfilesapi.Config) (int, error) {
	u := &nasfilesapi.User{}

	reqBuffer := new(bytes.Buffer)
	reqBuffer.ReadFrom(r.Body)

	// parse json
	err := json.Unmarshal(reqBuffer.Bytes(), u)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	// Hash password and assign creation time
	u.SetPassword(u.Password)
	u.Created = time.Now()

	err = c.Services.User.Add(u)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	color.HiGreen("Created a new user with username %s...", u.Username)

	return jsonPrint(w, u)
}

//UserGet retrieves a user from the database and returns it
func UserGet(w http.ResponseWriter, r *http.Request, c *nasfilesapi.Config) (int, error) {
	uid := mux.Vars(r)["uid"]

	var u *nasfilesapi.User

	u, err := c.Services.User.GetByUsername(uid)
	if err != nil {
		return http.StatusNotFound, err
	}

	return utils.JSONPrint(w, u)
}

//UserDelete inserts a user into the database
func UserDelete(w http.ResponseWriter, r *http.Request, c *nasfilesapi.Config) (int, error) {
	uid := mux.Vars(r)["uid"]

	err := c.Services.User.Delete(uid)
	if err != nil {
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
