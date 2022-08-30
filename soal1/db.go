package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func DB() *sql.DB {
	db_root := os.Getenv("DB")
	db, err := sql.Open("mysql", db_root)

	PanicIfError(err)

	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(3 * time.Second)
	// db.SetConnMaxLifetime(68 * time.Minute)
	// db.SetConnMaxIdleTime(10 * time.Minute)

	errormessage := ""
	if err != nil {
		errormessage = "Got error in mysql connector : " + err.Error()
		fmt.Println(errormessage)
	}

	err = db.Ping()
	if db.Ping() != nil {
		errormessage = "Could not connect to database : " + err.Error()
		fmt.Println(errormessage)
	}

	return db
}
