package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {

	var domain = flag.String("d", "", "Domain to analyse")
	var filename = flag.String("o", "sitemap.xml", "Output filename")
	var simultaneus = flag.Int("s", 3, "Number of concurrent connections")
	flag.Parse()

	if *domain == "" {
		fmt.Println("URL can not be empty")
		os.Exit(1)
	}

	fmt.Println("Domain:", *domain)
	fmt.Println("Simultaneus:", *simultaneus)

	if *simultaneus < 1 {
		fmt.Println("There can't be less than 1 simulataneous conexions")
		os.Exit(1)
	}

	scanning := make(chan int, *simultaneus) // Semaphore
	newLinks := make(chan []Links, 10000)    // New links to scan
	pages := make(chan Page, 10000)          // Pages scanned
	started := make(chan int, 10000)         // Crawls started
	finished := make(chan int, 10000)        // Crawls finished

	var indexed, noIndex []string

	seen := make(map[string]bool)

	start := time.Now()

	defer func() {

		close(newLinks)
		close(pages)
		close(started)
		close(finished)
		close(scanning)

		fmt.Printf("\nTime finished sitemap %s\n", time.Since(start))
		fmt.Printf("Index: %6d - NoIndex %6d\n", len(indexed), len(noIndex))
	}()

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

			if page.NoIndex {

				if !isUrlInSlice(page.Url, noIndex) {
					noIndex = append(noIndex, page.Url)
				}

			} else {

				if !isUrlInSlice(page.Url, indexed) {
					indexed = append(indexed, page.Url)
				}
			}
		}

		// Break the for loop once all scans are finished
		if len(started) > 1 && len(scanning) == 0 && len(started) == len(finished) {
			break
		}
	}

	createSitemap(indexed, *filename)
}
