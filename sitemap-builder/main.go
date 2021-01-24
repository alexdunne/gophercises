package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type urlset struct {
	Urls []loc `xml:"url"`
}

type loc struct {
	Value string `xml:"loc"`
}

func main() {
	url := flag.String("url", "", "URL starting point")
	flag.Parse()

	if *url == "" {
		exit("A url is required")
	}

	// 1. Download the contents of a page
	// 2. Parse the contents to find all links in the page
	// 3. Download unvisited pages and repeat the above
	// Once all pages have been visited then we're done

	urls := fetchAllSiteURLs(*url, *url)

	var xmlToConvert urlset

	for _, url := range urls {
		xmlToConvert.Urls = append(xmlToConvert.Urls, loc{Value: url})
	}

	encoder := xml.NewEncoder(os.Stdout)
	encoder.Indent("", "  ")
	err := encoder.Encode(&xmlToConvert)
	if err != nil {
		panic(err)
	}

	fmt.Println()
}

func fetchAllSiteURLs(domain, startURL string) []string {
	visitedURLs := make(map[string]bool)
	urlsToVisit := []string{startURL}

	for len(urlsToVisit) > 0 {
		currentURL := urlsToVisit[0]
		visitedURLs[currentURL] = true

		contents, err := getPageContent(currentURL)
		if err != nil {
			fmt.Println(err)
			exit(fmt.Sprintf("Something went wrong whilst fetching the content for %s", currentURL))
		}

		links, err := Parse(contents)
		if err != nil {
			fmt.Println(err)
			exit(fmt.Sprintf("Something went wrong whilst extracting the links for %s", currentURL))
		}

		urls := filterDuplicateURLs(links)

		for _, url := range urls {

			// Different domain so we don't need to include it
			if !isSameDomain(domain, url) {
				continue
			}

			formattedURL := formatURL(domain, url)

			// Already visited the domain so skip
			if visitedURLs[formattedURL] {
				continue
			}

			urlsToVisit = append(urlsToVisit, formattedURL)
		}

		urlsToVisit = urlsToVisit[1:]
	}

	var ret []string

	for url := range visitedURLs {
		ret = append(ret, url)
	}

	return ret
}

func getPageContent(url string) (io.Reader, error) {
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func formatURL(domain, url string) string {
	if strings.HasPrefix(url, domain) {
		return url
	}

	return domain + url
}

func isSameDomain(domain, url string) bool {
	if strings.HasPrefix("/", url) {
		return true
	}

	return strings.HasPrefix(url, domain)
}

func filterDuplicateURLs(links []Link) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, link := range links {
		if _, ok := keys[link.Href]; !ok {
			keys[link.Href] = true
			list = append(list, link.Href)
		}
	}

	return list
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
