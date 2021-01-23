package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
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

func buildLink(n *html.Node) Link {
	var link Link

	for _, att := range n.Attr {
		if att.Key == "href" {
			link.Href = att.Val
		}
	}

	link.Text = getNodeText(n)

	return link
}

func getNodeText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.Type != html.ElementNode {
		return ""
	}

	var ret string

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += getNodeText(c)
	}

	return strings.Join(strings.Fields(ret), " ")
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
