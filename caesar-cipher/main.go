package main

import (
	"fmt"
)

func main() {
	var input string
	var cipherKey int

	fmt.Scanf("%s\n", &input)
	fmt.Scanf("%d\n", &cipherKey)

	alphabetLower := "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	output := ""

	for _, char := range input {
		switch {
		case findOriginalPosition(alphabetLower, char) != -1:
			pos := findOriginalPosition(alphabetLower, char)
			newPos := (pos + cipherKey) % len(alphabetLower)

			output += string(alphabetLower[newPos])
		case findOriginalPosition(alphabetUpper, char) != -1:
			pos := findOriginalPosition(alphabetUpper, char)
			newPos := (pos + cipherKey) % len(alphabetUpper)

			output += string(alphabetUpper[newPos])
		default:
			output += string(char)

		}
	}

	fmt.Println(output)
}

func findOriginalPosition(alphabet string, value rune) int {
	for i, letter := range alphabet {
		if letter == value {
			return i
		}
	}

	return -1
}
