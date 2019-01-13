package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	var domain = flag.String("u", "", "URL to extract")
	flag.Parse()

	if *domain == "" {
		fmt.Println("URL can not be empty")
		os.Exit(1)
	}

	links, err := getLinks(*domain)
	if err != nil {
		fmt.Println("Could not get any links from Startsite")
	}

	for _, link := range links {
		newLinks, err := getLinks(link.Href)
		if err != nil {
			break
		}

		for _, new := range newLinks {
			if !doesLinkExist(new, links) {
				links = append(links, new)
			}
		}
	}

	for _, val := range links {
		fmt.Println(val)
	}
}
