package main

import "net/http"

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := app.models.Token.AuthenticateToken(r)
		if err != nil{
			payload := jsonResponse{
				Error: true,
				Message: "invalid auth credentials",				
			}

			_ = app.writeJSON(w, http.StatusUnauthorized, payload)
			return
		}
		// If we pass the error check
		next.ServeHTTP(w, r)
	})
}