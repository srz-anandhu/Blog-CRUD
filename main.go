package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	user     = "postgres"
	password = "password"
	host     = "localhost"
	port     = 5432
	dbname   = "BlogDB"
)

var Db *sql.DB

func InitDB() {
	connectionString := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", user, password, host, port, dbname)
	Db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("check the connection string: %v", err)
	}

	defer Db.Close()

	if err := Db.Ping(); err != nil {
		log.Fatalf("unable to connect to the database : %v", err)
	}

	fmt.Println("successfully connected to db...")
}

func main() {
	InitDB()
	defer Db.Close()
}
