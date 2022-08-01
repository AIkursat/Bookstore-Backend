package main

import (
	"Bookstore-Backend/internal/data"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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
	
	// we have a valid user, let's generate a toke

	// If user active and valid, generate a token

   if user.Active == 0 {
	app.errorJSON(w, errors.New("user is not active"))
	return
   }

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
func (app *application) AllUsers(w http.ResponseWriter, r *http.Request){
	var users data.User
	all, err := users.GetAll()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: "success",
		Data: envelope{"users": all},
	}
	app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) EditUser(w http.ResponseWriter, r *http.Request){
	var user data.User
	err := app.readJSON(w, r, &user)
	if err != nil{
		app.errorJSON(w, err)
		return
	}

	// what is the id of the in my json payload?

	if user.ID == 0 { // 0 means it doesn't exist, then add.
		// add user
		if _, err := app.models.User.Insert(user); 
		 err != nil {
			app.errorJSON(w, err)
		return
		}
	} else{
		// Edit user, u is for getting the user
		u, err := app.models.User.GetOne(user.ID)
		if err != nil{
			app.errorJSON(w, err)
		    return
		}

		u.Email = user.Email
		u.FirstName = user.FirstName
		u.LastName = user.LastName
		u.Active = user.Active

		if err := u.Update();
		 err != nil {
			app.errorJSON(w, err)
		    return
		 }

		 // if password != string, update password
		 if user.Password != "" {
			err := u.ResetPassword(user.Password)
			if err != nil{
				app.errorJSON(w, err)
				return
			}
		 }
	}

    payload := jsonResponse{
		Error: false,
		Message: "Changes saved",
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *application) Getuser(w http.ResponseWriter, r *http.Request){
	userId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user, err := app.models.User.GetOne(userId)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, user)
}

func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request){
   var requestPaylaod struct {
	ID int `json:"id"`
   }

   err := app.readJSON(w, r, &requestPaylaod)
   if err != nil{
	app.errorJSON(w, err)
	return
   }

   err = app.models.User.DeleteById(requestPaylaod.ID)
   if err != nil{
	app.errorJSON(w, err)
	return
   }

   payload := jsonResponse{
	Error: false,
	Message: "User deleted",
   }

   _ = app.writeJSON(w, http.StatusOK, payload) // ignored error

}

func (app *application) LogUserOutAndSetInactive(w http.ResponseWriter, r *http.Request){
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user, err := app.models.User.GetOne(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user.Active = 0
	err = user.Update()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// delete tokens for user
	err =  app.models.Token.DeleteTokensForuser(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: "user logged out and set to inactive",
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *application) ValidateToken(w http.ResponseWriter, r *http.Request){
    var requestPaylaod struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &requestPaylaod)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
    
	valid := false
	valid, _ = app.models.Token.ValidToken(requestPaylaod.Token)
	
	payload := jsonResponse {
		Error: false,
		Data: valid,
	}
    
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) AllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := app.models.Book.GetAll()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse {
		Error: false,
		Message: "success",
		Data: envelope{"books": books},
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) OneBook(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	book, err := app.models.Book.GetOneBySlug(slug)
	if err != nil{
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse {
		Error: false,
		Data: book,
	}
	app.writeJSON(w, http.StatusOK, payload)
} 
func (app *application) AuhtorsAll(w http.ResponseWriter, r *http.Request) {
	all, err := app.models.Author.All()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	type selectData struct {
		Value int `json:"value"`
		Text string `json:"text"`
	}

	var results []selectData

	for _, x := range all {
		author := selectData {
			Value: x.ID,
			Text: x.AuthorName,
		}

		results = append(results, author)
	}

	payload := jsonResponse{
		Error: false,
		Data: results,
	}

	app.writeJSON(w, http.StatusOK, payload)
}