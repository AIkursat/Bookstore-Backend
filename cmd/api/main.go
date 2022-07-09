package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type config struct {
	port int
}

type application struct {
	config  config // sharing configiration with application
	infoLog *log.Logger // Logger
	errorLog *log.Logger
}


func main() {
    var cfg config 
	cfg.port = 8081 // will use 8081 port

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		config: cfg,
		infoLog: infoLog,
		errorLog: errorLog,
	}

	err := app.serve()
	if err != nil{
       log.Fatal(err)
	}

}

func(app *application) serve() error{

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		var payload struct{
			Okay bool `json:"okay"`
			Message string `json:"message"`
		}
		payload.Okay = true
		payload.Message = "Hello, World"

		out, err := json.MarshalIndent(payload, "", "\t") // Make json easy to read to payload, no prefix "", tab
		
		if err != nil{
            app.errorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) //200
		w.Write(out)

	})

	app.infoLog.Println("API Listening on port", app.config.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", app.config.port), nil) // d: decimal integer

}  