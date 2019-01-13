package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {

	var domain = flag.String("u", "", "URL to extract")
	var filename = flag.String("s", "sitemap.xml", "Filename for the sitemap")
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

	scanned := []string{}

	// Extract Links from Startsite
	links, err := getLinks(*domain, root)
	if err != nil {
		fmt.Println("Could not get any links from Startsite")
	}

	for i := 0; i < 10; i++ {

		fmt.Printf("Round %d Links found: %d\n", i, len(links))

		for _, link := range links {

			if isLinkScanned(link.Href, scanned) || link.NoFollow {
				continue
			}

			newLinks, err := getLinks(link.Href, root)
			if err != nil {
				break
			}

			scanned = append(scanned, link.Href)

			for _, new := range newLinks {
				if !doesLinkExist(new, links) {
					links = append(links, new)
				}
			}
		}

	}

	fmt.Printf("Total Links found: %d\n", len(links))

	createSitemap(links, *filename)
}
