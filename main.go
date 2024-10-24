package main

import (
	"fmt"
	"sync"

	"github.com/pratikdev/url-scrapper/utils"
)

func main() {
	var wg = &sync.WaitGroup{}
	var mut = &sync.Mutex{}

	results := []string{}
	var seen = []string{}

	var URL string

	fmt.Println("Enter URL to scrape:")
	fmt.Scanln(&URL)

	wg.Add(1)
	go utils.GetValidURLs(wg, mut, &results, &seen, URL)
	wg.Wait()

	fmt.Println("") // line break
	fmt.Println("Total URLs found:", len(results))
	fmt.Println("Total seen URLs:", len(seen))

	fmt.Println("") // line break
	fmt.Println("Valid URLs:")
	for _, url := range results {
		fmt.Println(url)
	}
}
