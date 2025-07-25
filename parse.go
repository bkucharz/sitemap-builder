package link

import (
	"errors"
	"fmt"
	"io"
	"os"
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
		fmt.Printf("Cannot parse HTML file %v\n", r)
		os.Exit(1)
	}
	linkNodes := filterLinkNodes(doc)
	links := getLinks(linkNodes)
	return links, nil
}

func filterLinkNodes(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}
	var ret []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, filterLinkNodes(c)...)
	}
	return ret
}

func getLink(node *html.Node) (Link, error) {
	link := Link{}
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
			break
		}
	}
	if link.Href == "" {
		return Link{}, errors.New("node missing href attribute")
	}
	text := getText(node)
	link.Text = strings.Join(strings.Fields(text), " ")
	return link, nil
}

func getText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	var text string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		text += getText(c)
	}
	return text
}

func getLinks(nodes []*html.Node) []Link {
	var links []Link
	for _, node := range nodes {
		link, err := getLink(node)
		if err != nil {
			continue
		}
		links = append(links, link)
	}
	return links
}
