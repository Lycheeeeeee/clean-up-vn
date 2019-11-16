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
	usersdata := models.GetAllUsers()
	resp := u.Message(true, "success")
	resp["usersdata"] = usersdata
	u.Respond(w, resp)
}

var UpdateUser = func(w http.ResponseWriter, r *http.Request) {
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

var DownloadByte = func(w http.ResponseWriter, r *http.Request){
	





}
