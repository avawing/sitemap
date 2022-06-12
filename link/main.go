package link

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

// Link represents a link in an html document
type Link struct {
	Href string
	Text string
}

// Parse takes in an html document, returns slice of links parsed from document
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	nodes := linkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	//dfs(doc, "")
	// unneeded dfs now
	return links, nil
}

func parseText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	if node.Type != html.ElementNode {
		return ""
	}
	var ret string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		ret += parseText(child) + " "
	}

	return strings.Join(strings.Fields(ret), " ")
}

func buildLink(node *html.Node) Link {
	var ret Link

	for _, attr := range node.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = parseText(node)
	return ret
}

func linkNodes(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}
	var ret []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		ret = append(ret, linkNodes(child)...)
	}

	return ret
}
