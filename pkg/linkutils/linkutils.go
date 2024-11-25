package linkutils

import (
	"net/url"
	"strings"
)

type LinkUtils struct{}

func NewLinkUtils() *LinkUtils {
	return &LinkUtils{}
}

func (l *LinkUtils) ToAbsoluteURL(href string, base string) string {
	parsedBase, err := url.Parse(base)
	if err != nil {
		return href
	}
	parsedHref, err := url.Parse(href)
	if err != nil {
		return href
	}
	return parsedBase.ResolveReference(parsedHref).String()
}

func (l *LinkUtils) IsInternalLink(link string, base string) bool {
	parsedBase, err := url.Parse(base)
	if err != nil {
		return false
	}
	parsedLink, err := url.Parse(link)
	if err != nil {
		return false
	}
	return strings.EqualFold(parsedBase.Host, parsedLink.Host)
}
