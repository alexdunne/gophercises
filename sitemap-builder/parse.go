package main

import (
	"io"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)

	if err != nil {
		return nil, err
	}

	nodes := findLinkNodes(doc)

	var links []Link

	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil
}

func findLinkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var linkNodes []*html.Node

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		linkNodes = append(linkNodes, findLinkNodes(c)...)
	}

	return linkNodes
}

func buildLink(n *html.Node) Link {
	var link Link

	for _, att := range n.Attr {
		if att.Key == "href" {
			link.Href = att.Val
		}
	}

	return link
}
