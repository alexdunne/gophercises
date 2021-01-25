package main

import (
	"regexp"
)

// Normalise removes all non-number characters from a string
func Normalise(phoneNumber string) string {
	regex := regexp.MustCompile("\\D")
	return regex.ReplaceAllString(phoneNumber, "")
}
