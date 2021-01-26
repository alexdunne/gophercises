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
		fmt.Printf("Normalising %s (id: %d)\n", phoneNumber.value, phoneNumber.id)

		normalisedNumber := normalise(phoneNumber.value)

		if phoneNumber.value == normalisedNumber {
			fmt.Println("Nothing to see here")
			continue
		}

		fmt.Printf("Found a difference (before: '%s', after: '%s')\n", phoneNumber.value, normalisedNumber)

		existing, err := getPhoneByValue(db, normalisedNumber)
		must(err)

		// Value doesn't exist in the DB so update the current record
		if existing == nil {
			phoneNumber.value = normalisedNumber
			must(updatePhone(db, phoneNumber))
			continue
		}

		// Value already exists at least once so delete the current record
		must(deletePhone(db, phoneNumber))
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

func getPhoneByValue(db *sql.DB, value string) (*phoneNumber, error) {
	var p phoneNumber

	err := db.QueryRow("SELECT id, value from phone_numbers WHERE value = $1", value).Scan(&p.id, &p.value)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &p, err
}

func updatePhone(db *sql.DB, phone phoneNumber) error {
	_, err := db.Exec("UPDATE phone_numbers SET value = $2 WHERE id = $1", phone.id, phone.value)
	return err
}

func deletePhone(db *sql.DB, phone phoneNumber) error {
	_, err := db.Exec("DELETE FROM phone_numbers WHERE id = $1", phone.id)
	return err
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
