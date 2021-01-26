package main

import (
	"fmt"
	"regexp"

	"normaliser/db"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "normaliser"
	password = "normaliser"
	dbname   = "normaliser"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)

	must(db.Migrate("postgres", psqlInfo))

	db, err := db.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	must(db.Seed())

	phoneNumbers, err := db.GetAllPhoneNumbers()
	must(err)

	for _, phoneNumber := range phoneNumbers {
		fmt.Printf("Normalising %s (id: %d)\n", phoneNumber.Value, phoneNumber.ID)

		normalisedNumber := normalise(phoneNumber.Value)

		if phoneNumber.Value == normalisedNumber {
			fmt.Println("Nothing to see here")
			continue
		}

		fmt.Printf("Found a difference (before: '%s', after: '%s')\n", phoneNumber.Value, normalisedNumber)

		existing, err := db.GetPhoneByValue(normalisedNumber)
		must(err)

		// Value doesn't exist in the DB so update the current record
		if existing == nil {
			phoneNumber.Value = normalisedNumber
			must(db.UpdatePhone(&phoneNumber))
			continue
		}

		// Value already exists at least once so delete the current record
		must(db.DeletePhone(phoneNumber.ID))
	}
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
