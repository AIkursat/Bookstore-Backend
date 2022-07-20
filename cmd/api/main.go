package main

import (
	"Bookstore-Backend/internal/data"
	"Bookstore-Backend/internal/driver"
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
	models data.Models
}


func main() {
    var cfg config 
	cfg.port = 8081 // will use 8081 port

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
   
	// dsn means Data Source Name
	dsn := "host=localhost port=5432 user=postgres password=0123321 dbname=bookkeeper sslmode=disable timezone=UTC connect_timeout=5"
    db, err := driver.ConnectPostgres(dsn) 
	if err != nil{
		log.Fatal("cannot connect to database")
	}

	defer db.SQL.Close()


	app := &application{
		config: cfg,
		infoLog: infoLog,
		errorLog: errorLog,
		models: data.New(db.SQL),
	}

	err = app.serve()
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