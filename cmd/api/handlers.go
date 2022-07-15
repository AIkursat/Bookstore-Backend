package main

import (
	"net/http"
)

// jsonResponse is the type used for generic JSON responses
type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

// Login is the handler used to attempt to log a user into the api
func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	type credentials struct {
		UserName string `json:"email"`
		Password string `json:"password"`
	}

	var creds credentials // It keeps a place for credentials
	var payload jsonResponse

	err := app.readJSON(w, r, &creds) // third one is what we wanna decode the json into
    if err != nil{
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = "invalid json supplied, or json missing entirely"
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
	}   
 
	// TODO authenticate
	app.infoLog.Println(creds.UserName, creds.Password)

	// send back a response
	payload.Error = false
	payload.Message = "Signed in"

	// out, err := json.MarshalIndent(payload, "", "\t")
	err = app.writeJSON(w, http.StatusOK, payload) // used instead of previos one
	if err != nil {
		app.errorLog.Println(err) 
	}
}

