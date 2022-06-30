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

	tram := strings.Split(link, "#")[0]

	if !*useQueries {
		tram = removeQuery(tram)
	}

	return tram
}

func isInternLink(link string) bool {
	return strings.Index(link, root) == 0
}

func removeQuery(link string) string {
	return strings.Split(link, "?")[0]
}

func isStart(link string) bool {
	return strings.Compare(link, root) == 0
}

func isValidExtension(link string) bool {
	for _, extension := range Extensions {
		if strings.Contains(strings.ToLower(link), extension) {
			return false
		}
	}
	return true
}

func isValidLink(link string) bool {

	if isInternLink(link) && !isStart(link) && isValidExtension(link) {
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

func isUrlInSlice(search string, array []string) bool {

	withSlash := search[:len(search)-1]
	withoutSlash := search

	if string(search[len(search)-1]) == "/" {
		withSlash = search
		withoutSlash = search[:len(search)-1]
	}

	for _, val := range array {
		if val == withSlash || val == withoutSlash {
			return true
		}
	}

	return false
}

func getLinks(domain string) (page Page, err error) {

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
						if isValidLink(foundLink) {
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

func takeLinks(toScan string, started chan int, finished chan int, scanning chan int, newLinks chan []Links, pages chan Page) {

	started <- 1
	scanning <- 1
	defer func() {
		<-scanning
		finished <- 1
		fmt.Printf("\rStarted: %6d - Finished %6d", len(started), len(finished))
	}()

	// Get links
	page, err := getLinks(toScan)
	if err != nil {
		return
	}

	// Save Page
	pages <- page

	// Save links
	newLinks <- page.Links
}

func createSitemap(links []string, filename string) {

	var total = []byte(xml.Header)
	total = append(
		total,
		[]byte(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)...,
	)
	total = append(total, []byte("\n")...)

	for _, val := range links {
		pos := UrlSitemap{Loc: val}
		output, err := xml.MarshalIndent(pos, "  ", "    ")
		check(err)
		total = append(total, output...)
		total = append(total, []byte("\n")...)
	}

	total = append(total, []byte(`</urlset>`)...)

	err := ioutil.WriteFile(filename, total, 0644)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
