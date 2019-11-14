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

var CreateProject = func(w http.ResponseWriter, r *http.Request) {

	pro := &models.Project{}
	// layout := "2006-01-02T15:04:05.000Z"
	// str := r.Body.time
	// t, err := time.Parse(layout, str)
	// if err != nil {
	// fmt.Println(err)
	// }
	// r.Body.time = t
	err := json.NewDecoder(r.Body).Decode(pro) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := pro.Create() //Create account
	u.Respond(w, resp)
}

var GetAllProjects = func(w http.ResponseWriter, r *http.Request) {
	projectsdata := models.GetAllProjects()
	resp := u.Message(true, "success")
	resp["projectsdata"] = projectsdata
	u.Respond(w, resp)
}

var GetProjectByID = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	projectdata := models.GetProject(params["id"])
	if projectdata == nil {
		resp := u.Message(false, "Invalid Id")
		u.Respond(w, resp)
	} else {
		resp := u.Message(true, "success")
		resp["projectdata"] = projectdata
		u.Respond(w, resp)
	}
}

var InputResult = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	project := &models.Project{}
	err := json.NewDecoder(r.Body).Decode(project) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	u64, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	project.ID = uint(u64)
	project.Status = "close"
	resp := project.InputResultNCloseProject()
	u.Respond(w, resp)
}
