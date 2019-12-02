package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Lycheeeeeee/clean-up-vn/models"
	u "github.com/Lycheeeeeee/clean-up-vn/utils"
	"github.com/gorilla/mux"
)

var CreateUser = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := user.Create() //Create account
	u.Respond(w, resp)
}

var GetUserByID = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	userdata := models.GetUser(params["id"])
	if userdata == nil {
		resp := u.Message(false, "Invalid Id")
		u.Respond(w, resp)
	} else {
		resp := u.Message(true, "success")
		resp["userdata"] = userdata
		u.Respond(w, resp)
	}
}

var GetAllUsers = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	usersdata := models.GetAllUsers()
	resp := u.Message(true, "success")
	resp["usersdata"] = usersdata
	u.Respond(w, resp)
}

var UpdateUser = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	u64, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	user.ID = uint(u64)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := user.UpdateSub()
	u.Respond(w, resp)
}

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := user.CreateAccount() //Create account
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(user.Email, user.Password)
	u.Respond(w, resp)
}

var Socialauthenticate = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	userdata := models.GetTokenFromSocial(params["id"])
	if userdata == nil {
		resp := u.Message(false, "Invalid Id")
		u.Respond(w, resp)
	} else {
		resp := u.Message(true, "success")
		resp["userdata"] = userdata
		u.Respond(w, resp)
	}
}
