package main

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "normaliser"
	password = "normaliser"
	dbname   = "normaliser"
)

func main() {
	db := createDBConnection()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		panic(err)
	}
}

func createDBConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	return db
}

// normalise removes all non-number characters from a string
func normalise(phoneNumber string) string {
	regex := regexp.MustCompile("\\D")
	return regex.ReplaceAllString(phoneNumber, "")
}
