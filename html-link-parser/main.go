package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	filename := flag.String("file", "", "The filename to parse")
	flag.Parse()

	if *filename == "" {
		exit(fmt.Sprintf("A file option must be provided"))
	}

	file, err := os.Open(*filename)

	if err != nil {
		exit(fmt.Sprintf("Something went wrong whilst opening the file %s", *filename))
	}

	doc, err := html.Parse(file)

	if err != nil {
		exit(fmt.Sprintf("Something went wrong whilst parsing the file %s", *filename))
	}

}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
