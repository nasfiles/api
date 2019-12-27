package http

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nasfiles/api"
	"github.com/nasfiles/api/utils"
)

//UserAdd inserts a user into the database
func UserAdd(w http.ResponseWriter, r *http.Request, c *api.Config) (int, error) {
	u := &api.User{
		UID:      base64.URLEncoding.EncodeToString(utils.GenerateSHA256()),
		Username: "fabiofcferreira",
		Email:    "ffcfpten@gmail.com",
		Name:     "Fábio",
		Created:  time.Now(),
	}

	c.Services.User.Add(u)

	return http.StatusOK, nil
}

//UserGet retrieves a user from the database and returns it
func UserGet(w http.ResponseWriter, r *http.Request, c *api.Config) (int, error) {
	uid := mux.Vars(r)["uid"]

	var u *api.User

	u, err := c.Services.User.GetByUsername(uid)
	if err != nil {
		return http.StatusNotFound, err
	}

	return utils.JSONPrint(w, u)
}

//UserDelete inserts a user into the database
func UserDelete(w http.ResponseWriter, r *http.Request, c *api.Config) (int, error) {
	uid := mux.Vars(r)["uid"]

	err := c.Services.User.Delete(uid)
	if err != nil {
		return http.StatusNotFound, err
	}

	return http.StatusOK, nil
}
