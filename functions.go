package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
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

func getLinks(domain string, root string) (page Page, err error) {

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

	page.Url = domain
	foundMeta := false

	var f func(*html.Node)
	f = func(n *html.Node) {

		if n.Type == html.ElementNode && n.Data == "meta" {
			for _, a := range n.Attr {
				if a.Key == "name" && a.Val == "robots" {
					foundMeta = true
				}
				if foundMeta {
					if a.Key == "content" && strings.Contains(a.Val, "noindex") {
						page.NoIndex = true
					}
					if a.Key == "content" && strings.Contains(a.Val, "nofollow") {
						page.NoFollow = true
					}
				}
			}
		}

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

			if ok && !doesLinkExist(newLink, page.Links) {
				page.Links = append(page.Links, newLink)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

	}
	f(doc)

	return
}

func takeLinks(domain string, root string, savedLinks chan string) {

}

func createSitemap(links []Links, filename string) {

	var total = []byte(xml.Header)
	total = appendBytes(total, []byte(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`))

	for _, val := range links {
		pos := UrlSitemap{Loc: val.Href}
		output, err := xml.Marshal(pos)
		check(err)

		for _, b := range output {
			total = append(total, b)
		}
	}

	total = appendBytes(total, []byte(`</urlset>`))

	err := ioutil.WriteFile(filename, total, 0644)
	check(err)
}

func appendBytes(appendTo []byte, toAppend []byte) []byte {
	for _, val := range toAppend {
		appendTo = append(appendTo, val)
	}

	return appendTo
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
