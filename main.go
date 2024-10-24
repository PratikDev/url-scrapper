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

	const URL = "https://scrape-me.dreamsofcode.io"

	wg.Add(1)
	go utils.GetValidURLs(wg, mut, &results, &seen, URL)
	wg.Wait()

	fmt.Println("Total URLs found:", len(results))
	fmt.Println("Total seen URLs:", len(seen))
}
