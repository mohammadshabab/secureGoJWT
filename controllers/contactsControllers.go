package controllers

import (
	"encoding/json"
	"net/http"
	"securegojwt/app"
	"securegojwt/models"
	"securegojwt/utils"
	"strconv"

	"github.com/gorilla/mux"
)

var CreateContact = func(w http.ResponseWriter, r *http.Request) {
	var usr app.ContextKey = "user"
	user := r.Context().Value(usr).(uint) //Get the id of user who send request
	contact := &models.Contact{}

	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error while decoding request body"))
		return
	}
	contact.UserId = user
	resp := contact.Create()
	utils.Respond(w, resp)
}

var GetContactsFor = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		utils.Respond(w, utils.Message(false, "Error in request"))
		return
	}
	data := models.GetContacts(uint(id))
	resp := utils.Message(true, "success")
	resp["data"] = data
	utils.Respond(w, resp)
}
