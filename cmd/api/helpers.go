package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error{
    
	maxBytes := 1048576 //  max file size will accept as post. It means 1 mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes)) // maxBytes converted to the int64


	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil{
		return nil
	}
	 
	err = dec.Decode(&struct{}{}) // to be make sure body has single json value

	if err != io.EOF{ // like nil
		return errors.New("body must have only a single json value")
	}

	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header ) error{

    // headers ...http.Header mean we can put one, more or none headers.
      
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil{
		return err
	}

	if len(headers) > 0 {
     for key, value := range headers[0]{
		 w.Header()[key] = value
	 }

	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil{
		return err
	}
  
	return nil
}

func(app *application) errorJSON(w http.ResponseWriter, err error, status ...int){
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var customErr error

	switch{
	case strings.Contains(err.Error(), "SQLSTATE 23505"): // Postgres error
	customErr = errors.New("duplicate value violates unique constraints")
	statusCode = http.StatusForbidden
	case strings.Contains(err.Error(), "SQLSTATE 22001"):
		customErr = errors.New("too large value")
		statusCode = http.StatusForbidden
	case strings.Contains(err.Error(), "SQLSTATE 23403"):
		customErr = errors.New("Foreign key violation")
		statusCode = http.StatusForbidden
	default:
		customErr = err
		
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = customErr.Error()

	app.writeJSON(w, statusCode, payload)
}