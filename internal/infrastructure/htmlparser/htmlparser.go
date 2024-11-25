package htmlparser

import (
	"strings"

	"golang.org/x/net/html"

	"github.com/ZetoOfficial/domain-scraper/internal/domain"
	"github.com/ZetoOfficial/domain-scraper/pkg/linkutils"
)

type HTMLParser struct {
	linkUtils *linkutils.LinkUtils
}

func NewHTMLParser() *HTMLParser {
	return &HTMLParser{
		linkUtils: linkutils.NewLinkUtils(),
	}
}

func (p *HTMLParser) ParseLinks(htmlContent string, baseURL string) ([]domain.Link, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}

	var links []domain.Link
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var href string
			var text string
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href = attr.Val
					break
				}
			}
			if href != "" {
				href = p.linkUtils.ToAbsoluteURL(href, baseURL)
				if p.linkUtils.IsInternalLink(href, baseURL) {
					text = getText(n)
					links = append(links, domain.Link{
						Href:   href,
						Text:   text,
						Source: baseURL,
					})
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links, nil
}

func getText(n *html.Node) string {
	var text string
	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.TextNode {
			text += node.Data
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return strings.TrimSpace(text)
}
