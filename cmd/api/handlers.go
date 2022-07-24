package main

import (
	"errors"
	"net/http"
	"time"
)

// jsonResponse is the type used for generic JSON responses
type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

type envelope map[string] interface{} // adding envolpe

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

    
	// look up the user by email

	user, err := app.models.User.GetByEmail(creds.UserName)
	if err != nil{
		app.errorJSON(w, errors.New("invalid username/password"))
		return
	}

	// validate the user's password

	validPassword, err := user.PasswordMatches(creds.Password)
	if err != nil || !validPassword {
		app.errorJSON(w, errors.New("invalid username/password"))
		return
	}
	
	// we have a valid user, let's generate a token

	token, err := app.models.Token.GenerateToken(user.ID, 24 *time.Hour) // 24 hours expiry
	if err != nil{
		app.errorJSON(w, err)
		return
	}

	// save it to the db

	err = app.models.Token.Insert(*token, *user)
	if err != nil{
		app.errorJSON(w, err)
		return
	}

	// send back a response
	payload = jsonResponse{
		Error: false,
		Message: "Logged in",
		Data: envelope{"token": token, "user": user 	},
	}

	// out, err := json.MarshalIndent(payload, "", "\t")
	err = app.writeJSON(w, http.StatusOK, payload) // used instead of previos one
	if err != nil {
		app.errorLog.Println(err) 
	}
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request){
	var requestPayload struct{
		Token string `json:"token"`
	}


	err := app.readJSON(w, r, &requestPayload)
	if err != nil{
		app.errorJSON(w, errors.New("invalid json"))
		return
	}

	// if we pass that

	err = app.models.Token.DeleteByToken(requestPayload.Token)
	if err != nil{
		app.errorJSON(w, errors.New("invalid json"))
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: "logged out",
	}

	_ = app.writeJSON(w, http.StatusOK, payload) // we ignored the error

}
