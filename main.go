package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {

	var domain = flag.String("u", "", "URL to extract")
	var filename = flag.String("o", "sitemap.xml", "Output filename")
	var depth = flag.Int("d", 10, "Depth levels of crawling")

	flag.Parse()

	if *domain == "" {
		fmt.Println("URL can not be empty")
		os.Exit(1)
	}

	// Do First call to domain
	resp, err := http.Get(*domain)
	if err != nil {
		fmt.Println("Domain could not be reached!")
		return
	}
	// Todo: get favourite version of URL here
	defer resp.Body.Close()

	root := resp.Request.URL.String()

	links := []Links{}
	scanned := []string{}
	noIndex := []string{}
	noFollow := []string{}

	// Extract Links from Startsite
	page, err := getLinks(*domain, root)
	if err != nil {
		fmt.Println("Could not get any links from Startsite")
	}

	links = page.Links

	for i := 0; i < *depth; i++ {

		fmt.Printf("Round %d Links found: %d\n", i, len(links))

		for _, link := range links {

			if isLinkScanned(link.Href, scanned) || link.NoFollow {
				continue
			}

			scanned = append(scanned, link.Href)

			newLinks, err := getLinks(link.Href, root)
			if err != nil {
				break
			}

			if newLinks.NoFollow {
				noFollow = append(noFollow, newLinks.Url)
				continue
			}

			if newLinks.NoIndex {
				noIndex = append(noIndex, newLinks.Url)
			}

			for _, new := range newLinks.Links {
				if !doesLinkExist(new, links) {
					links = append(links, new)
				}
			}
		}

	}

	fmt.Printf("Total Links found: %d\n", len(links))

	// Remove noIndex
	if len(noIndex) > 0 {
		for _, no := range noIndex {
			for i, sub := range links {
				if strings.Compare(sub.Href, no) == 0 {
					links = append(links[:i], links[i+1:]...)
				}
			}
		}
	}

	createSitemap(links, *filename)
}
