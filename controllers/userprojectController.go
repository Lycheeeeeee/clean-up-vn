package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Lycheeeeeee/clean-up-vn/models"
	u "github.com/Lycheeeeeee/clean-up-vn/utils"
)

var CreateUserProject = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	userproject := &models.UserProject{}
	err := json.NewDecoder(r.Body).Decode(userproject)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := userproject.CreateUserProject()
	u.Respond(w, resp)
}

var DeleteUserProject = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	userproject := &models.UserProject{}
	err := json.NewDecoder(r.Body).Decode(userproject)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := userproject.LeaveProject()
	u.Respond(w, resp)
}
var Report = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resp := models.RunReport()
	u.Respond(w, resp)
}
