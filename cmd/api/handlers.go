package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {

	type credentials struct {
		UserName string `json:"email"`
		Password string `json:"password"`
	}

	var creds credentials
	var payload jsonResponse

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// Error message
		app.errorLog.Println("invalid json")
		payload.Error = true
		payload.Message = "invalid json"

		out, err := json.MarshalIndent(payload, "", "\t") // Make json easy to read to payload, no prefix "", tab

		if err != nil {
			app.errorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) //200
		w.Write(out)
		return
	}
		// Authhenticate

		app.infoLog.Println(creds.UserName, creds.Password)

        // Response

		payload.Error = false
		payload.Message = "Signed in"

		out, err := json.MarshalIndent(payload, "", "\t")

		if err != nil {
			app.errorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) //200
		w.Write(out)
}