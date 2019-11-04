// Copyright 2019 Antonio Sanchez (asanchez.dev). All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

var domain, root string
var filename *string
var simultaneus *int
var useQueries *bool

func main() {

	if len(os.Args) == 1 {
		fmt.Println("URL can not be empty")
		os.Exit(1)
	}
	domain = os.Args[1]

	filename = flag.String("o", "sitemap.xml", "Output filename")
	simultaneus = flag.Int("s", 3, "Number of concurrent connections")
	useQueries = flag.Bool("q", false, "Ignore queries on URLs")
	flag.Parse()

	fmt.Println("Domain:", domain)
	fmt.Println("Simultaneus:", *simultaneus)
	fmt.Println("Use Queries:", *useQueries)

	if *simultaneus < 1 {
		fmt.Println("There can't be less than 1 simulataneous conexions")
		os.Exit(1)
	}

	scanning := make(chan int, *simultaneus) // Semaphore
	newLinks := make(chan []Links, 100000)   // New links to scan
	pages := make(chan Page, 100000)         // Pages scanned
	started := make(chan int, 100000)        // Crawls started
	finished := make(chan int, 100000)       // Crawls finished

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
	resp, err := http.Get(domain)
	if err != nil {
		fmt.Println("Domain could not be reached!")
		return
	}
	// Todo: get favourite version of URL here
	defer resp.Body.Close()

	// Detected root domain
	root = resp.Request.URL.String()

	// Take the links from the startsite
	takeLinks(domain, started, finished, scanning, newLinks, pages)
	seen[domain] = true

	for {
		select {
		case links := <-newLinks:

			for _, link := range links {
				if !link.NoFollow {
					if !seen[link.Href] {
						seen[link.Href] = true
						go takeLinks(link.Href, started, finished, scanning, newLinks, pages)
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
