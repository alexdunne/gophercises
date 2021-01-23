package main

import (
	"flag"
	"fmt"
	"link"
	"net/http"
	"os"
)

func main() {
	filename := flag.String("filename", "", "The filename of the file to parse")
	url := flag.String("url", "", "A URL to parse")
	flag.Parse()

	if *filename == "" && *url == "" {
		exit("A filename or url is required")
	}

	var links []link.Link

	if *filename != "" {
		links = append(links, parseFile(*filename)...)
	}

	if *url != "" {
		links = append(links, parseURL(*url)...)
	}

	printLinks(links)
}

func parseFile(filename string) []link.Link {
	file, err := os.Open(filename)

	if err != nil {
		exit(fmt.Sprintf("The file %s could not be opened", filename))
	}

	links, err := link.Parse(file)

	if err != nil {
		exit(fmt.Sprintf("Something went wrong whilst parsing the %s file", filename))
	}

	return links
}

func parseURL(url string) []link.Link {
	response, err := http.Get(url)

	if err != nil {
		exit(fmt.Sprintf("Something went wrong whilst fetching the contents of the url %s", url))
	}

	links, err := link.Parse(response.Body)

	if err != nil {
		exit(fmt.Sprintf("Something went wrong whilst parsing the contents of the url %s", url))
	}

	return links
}

func printLinks(links []link.Link) {
	for _, line := range links {
		fmt.Printf("%+v\n", line)
	}
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
