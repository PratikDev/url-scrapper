package main

import (
	"fmt"
	"net/http"
	"slices"
	"sync"

	"github.com/pratikdev/url-scrapper/utils"
)

var wg = &sync.WaitGroup{}
var mut = sync.Mutex{}

var results = []string{}
var seen = []string{}

func main() {
	const URL = "https://scrape-me.dreamsofcode.io"

	wg.Add(1)
	go getValidURLs(URL)
	wg.Wait()

	fmt.Println("Total URLs found:", len(results))
	fmt.Println("Total seen URLs:", len(seen))
}

func getValidURLs(url string) {
	defer wg.Done()

	// if url is already seen, return
	if slices.Contains(seen, url) {
		return
	}

	statusCode, node := utils.FetchURL(url)

	mut.Lock()
	seen = append(seen, url) // add url to seen
	mut.Unlock()

	// if the status is not ok, return
	if statusCode != http.StatusOK {
		return
	}

	mut.Lock()
	results = append(results, url)
	mut.Unlock()

	hrefs := utils.GetHrefs(url, node)
	for _, href := range hrefs {
		wg.Add(1)
		// recursively call getValidURLs
		go getValidURLs(href)
	}
}
