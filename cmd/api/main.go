package main

import (
	
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

func(app *application) serve() error{ // serve was created was us

	app.infoLog.Println("API Listening on port", app.config.port)
    
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.config.port), // %d is the decimal, Addr is come from the Server
		Handler: app.routes(),	
	}

	return srv.ListenAndServe()

}  