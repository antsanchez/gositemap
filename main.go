package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {

	start := time.Now()

	newLinks := make(chan []Links, 10000)
	pages := make(chan Page, 10000)
	scanning := make(chan int, 5)
	started := make(chan int, 10000)
	finished := make(chan int, 10000)

	indexed := []Links{}
	seen := make(map[string]bool)

	defer func() {
		close(newLinks)
		close(pages)
		close(started)
		close(finished)
		close(scanning)
	}()

	var domain = flag.String("u", "", "URL to extract")
	var filename = flag.String("o", "sitemap.xml", "Output filename")
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

	// Detected root domain
	root := resp.Request.URL.String()

	// Take the links from the startsite
	takeLinks(*domain, root, started, finished, scanning, newLinks, pages)
	seen[*domain] = true

	for {
		select {
		case links := <-newLinks:
			for _, link := range links {
				if !link.NoFollow {
					if !seen[link.Href] {
						seen[link.Href] = true
						go takeLinks(link.Href, root, started, finished, scanning, newLinks, pages)
					}
				}
			}
		case page := <-pages:
			if !page.NoIndex {
				for _, link := range page.Links {
					indexed = append(indexed, link)
				}
			}
		}
		if len(started) > 1 && len(scanning) == 0 && len(started) == len(finished) {
			fmt.Printf("Breaking. Started: %d - Finished %d\n", len(started), len(finished))
			break
		}
	}

	fmt.Printf("Time finished crawling %s\n", time.Since(start))

	createSitemap(indexed, *filename)

	fmt.Printf("Time finished sitemap %s\n", time.Since(start))
}
