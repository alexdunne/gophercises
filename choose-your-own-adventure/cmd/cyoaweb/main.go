package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"cyoa"
)

func main() {
	filename := flag.String("file", "gopher.json", "The story JSON file")
	flag.Parse()
	fmt.Printf("Story %s loaded\n", *filename)

	file, err := os.Open(*filename)

	if err != nil {
		exit(fmt.Sprintf("Something went wrong whilst opening the file %s", *filename))
	}

	decoder := json.NewDecoder(file)
	var story cyoa.Story
	if err := decoder.Decode(&story); err != nil {
		exit(fmt.Sprintf("Something went wrong whilst decoding the file %s", *filename))
	}

	fmt.Printf("%+v\n", story)

}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
