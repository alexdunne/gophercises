package main

import (
	"context"
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

	must(migrate(db))
	must(seed(db))

	phoneNumbers, err := getAllPhoneNumbers(db)
	must(err)

	for _, phoneNumber := range phoneNumbers {
		fmt.Printf("%d: %s\n", phoneNumber.id, phoneNumber.value)
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

func migrate(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS phone_numbers`)
	if err != nil {
		return err
	}

	stmt := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)
	`
	_, err = db.Exec(stmt)

	return err
}

func seed(db *sql.DB) error {
	phoneNumbers := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}

	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, number := range phoneNumbers {
		_, err := tx.Exec("INSERT INTO phone_numbers(value) VALUES($1)", number)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

type phoneNumber struct {
	id    int
	value string
}

func getAllPhoneNumbers(db *sql.DB) ([]phoneNumber, error) {
	rows, err := db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []phoneNumber
	for rows.Next() {
		var p phoneNumber

		if err := rows.Scan(&p.id, &p.value); err != nil {
			return nil, err
		}

		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// normalise removes all non-number characters from a string
func normalise(phoneNumber string) string {
	regex := regexp.MustCompile("\\D")
	return regex.ReplaceAllString(phoneNumber, "")
}
