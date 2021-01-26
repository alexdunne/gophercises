package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

func Migrate(driverName string, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}

	_, err = db.Exec(`DROP TABLE IF EXISTS phone_numbers`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)
	`)
	if err != nil {
		return err
	}

	return db.Close()
}

func Open(driverName string, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

type DB struct {
	db *sql.DB
}

type PhoneNumber struct {
	ID    int
	Value string
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Seed() error {
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

	tx, err := db.db.BeginTx(ctx, nil)
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

func (db *DB) GetAllPhoneNumbers() ([]PhoneNumber, error) {
	rows, err := db.db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []PhoneNumber
	for rows.Next() {
		var p PhoneNumber

		if err := rows.Scan(&p.ID, &p.Value); err != nil {
			return nil, err
		}

		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (db *DB) GetPhoneByValue(value string) (*PhoneNumber, error) {
	var p PhoneNumber

	err := db.db.QueryRow("SELECT id, value from phone_numbers WHERE value = $1", value).Scan(&p.ID, &p.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &p, err
}

func (db *DB) UpdatePhone(phone *PhoneNumber) error {
	_, err := db.db.Exec("UPDATE phone_numbers SET value = $2 WHERE id = $1", phone.ID, phone.Value)
	return err
}

func (db *DB) DeletePhone(id int) error {
	_, err := db.db.Exec("DELETE FROM phone_numbers WHERE id = $1", id)
	return err
}
