package main

import (
	"fmt"
)

func main() {
	var input string
	var cipherKey int

	fmt.Scanf("%s\n", &input)
	fmt.Scanf("%d\n", &cipherKey)

	output := ""

	for _, char := range input {
		switch {
		case char >= 'A' && char <= 'Z':
			output += string(rotate('A', char, cipherKey))
		case char >= 'a' && char <= 'z':
			output += string(rotate('a', char, cipherKey))
		default:
			output += string(char)
		}
	}

	fmt.Println(output)
}

func rotate(base int, r rune, cipherKey int) rune {
	tmp := int(r) - base
	tmp = (tmp + cipherKey) % 26
	return rune(tmp + base)
}
