package driver

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

type DB struct{
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 5 // number of db connection
const maxIdleDbConn = 5 // When the task is complete the connection is marked as idle
const maxDbLifeTime = 5 * time.Minute // How long should they stay open


func ConnectPostgres(dsn string) (*DB, error) { // dsn = data source name, it'll return db and potentially error
	d, err := sql.Open("pgx", dsn) // dsn is connection string up here
	if err != nil{
		return nil, err
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifeTime)

	err = testDB(err, d)

	dbConn.SQL = d
	return dbConn, err
}  

func testDB(err error, d *sql.DB) error{
	err = d.Ping()
	if err != nil{
		fmt.Print("Error!", err)
	} else{
		fmt.Println("*** Pinged database succesfully ***")
	}
	return err
}
