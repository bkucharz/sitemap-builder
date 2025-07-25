package link

import (
	"fmt"
	"io"
	"os"

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
	dfs(doc, "")
	return nil, nil
}

func dfs(node *html.Node, padding string) {
	msg := node.Data
	if node.Type == html.ElementNode {
		msg = "<" + msg + ">"
	}
	fmt.Println(padding, msg)

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		dfs(child, padding+"  ")
	}
}
