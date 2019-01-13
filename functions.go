package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func SanitizeUrl(link string) string {

	for _, fal := range FalseUrls {
		if strings.Contains(link, fal) {
			return ""
		}
	}

	link = strings.TrimSpace(link)

	if string(link[len(link)-1]) != "/" {
		link = link + "/"
	}

	tram := strings.Split(link, "#")
	return tram[0]
}

func isInternLink(link string, root string) bool {

	if strings.Index(link, root) == 0 {
		return true
	}

	return false
}

func isStart(link string, root string) bool {

	if strings.Compare(link, root) == 0 {
		return true
	}

	return false
}

func isValidExtension(link string) bool {
	for _, extension := range Extensions {
		if strings.Contains(strings.ToLower(link), extension) {
			return false
		}
	}
	return true
}

func isValidLink(link string, root string) bool {

	if isInternLink(link, root) && !isStart(link, root) && isValidExtension(link) {
		return true
	}

	return false
}

func doesLinkExist(newLink Links, existingLinks []Links) (exists bool) {

	for _, val := range existingLinks {
		if strings.Compare(newLink.Href, val.Href) == 0 {
			exists = true
		}
	}

	return
}

func isLinkScanned(link string, scanned []string) (exists bool) {

	for _, val := range scanned {
		if strings.Compare(link, val) == 0 {
			exists = true
		}
	}

	return
}

func getLinks(domain string, root string) (links []Links, err error) {

	resp, err := http.Get(domain)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {

			ok := false
			newLink := Links{}

			for _, a := range n.Attr {

				if a.Key == "href" {
					link, err := resp.Request.URL.Parse(a.Val)
					if err == nil {
						foundLink := SanitizeUrl(link.String())
						if isValidLink(foundLink, root) {
							ok = true
							newLink.Href = foundLink
						}
					}
				}

				if a.Key == "rel" {
					if a.Val == "nofollow" {
						newLink.NoFollow = true
					}
				}

			}

			if ok && !doesLinkExist(newLink, links) {
				links = append(links, newLink)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return
}
