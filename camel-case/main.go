package main

import (
	"fmt"
	"unicode"
)

func main() {
	var phrase string

	fmt.Scanf("%s\n", &phrase)

	wordCount := 1

	for _, char := range phrase {
		if unicode.IsUpper(char) {
			wordCount++
		}
	}

	fmt.Println(wordCount)
}
