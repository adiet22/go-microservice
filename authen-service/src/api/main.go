package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/adiet95/go-microservice/authen-service/src/api/database"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var count int64

type Config struct {
	DB     *sql.DB
	Models database.Models
}

func main() {
	log.Println("Starting authen service")

	//TODO Connect to database
	conn := connectDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres !!")
	}

	//Set up config
	app := Config{
		DB:     conn,
		Models: database.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready")
			count++
		} else {
			log.Println("Connected to Postgres !!")
			return connection
		}
		if count > 10 {
			log.Println(err)
			return nil
		}
	}
}
