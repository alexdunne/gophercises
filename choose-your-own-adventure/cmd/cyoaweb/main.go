package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"cyoa"
)

func main() {
	port := flag.Int("port", 3000, "Web default port")
	filename := flag.String("file", "gopher.json", "The story JSON file")
	flag.Parse()
	fmt.Printf("Story %s loaded\n", *filename)

	file, err := os.Open(*filename)

	if err != nil {
		exit(fmt.Sprintf("Something went wrong whilst opening the file %s", *filename))
	}

	story, err := cyoa.JsonStory(file)
	if err != nil {
		exit(fmt.Sprintf("Something went wrong whilst decoding the file %s", *filename))
	}

	handler := cyoa.NewHandler(story)
	fmt.Printf("Starting port on :%d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
