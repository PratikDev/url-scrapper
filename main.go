package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var wg = &sync.WaitGroup{}
var mut = sync.Mutex{}

var results = []string{}

func checkURL(url string) bool {
	mut.Lock()
	inArray := slices.Contains(results, strings.TrimSuffix(url, "/"))
	mut.Unlock()
	return inArray
}

func main() {
	start := time.Now()

	const URL = "https://scrape-me.dreamsofcode.io"

	wg.Add(1)
	go getValidURLs(URL)
	wg.Wait()

	for _, result := range results {
		fmt.Println(result)
	}

	fmt.Println("Total URLs found:", len(results))

	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}

func getValidURLs(url string) {
	defer wg.Done()

	if checkURL(url) {
		return
	}

	statusCode, node := fetchURL(url)
	if statusCode == http.StatusOK {
		mut.Lock()
		results = append(results, url)
		mut.Unlock()

		hrefs := getHrefs(url, node)
		for _, href := range hrefs {
			wg.Add(1)
			go getValidURLs(href)
		}
	}
}

func fetchURL(url string) (statusCode int, node *html.Node) {
	fmt.Println("Checking URL:", url)

	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return response.StatusCode, nil
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	rawHtml := string(bytes)
	_node, err := html.Parse(strings.NewReader(rawHtml))
	if err != nil {
		panic(err)
	}

	return response.StatusCode, _node
}

func getHrefs(baseHost string, n *html.Node) []string {
	result := []string{}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" && isBaseHost(baseHost, a.Val) {
				parsedURL, err := url.Parse(baseHost)
				if err != nil {
					panic(err)
				}

				result = append(result, addPath(parsedURL.Scheme+"://"+parsedURL.Host, a.Val))
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result = append(result, getHrefs(baseHost, c)...)
	}
	return result
}

func isBaseHost(baseHost string, href string) bool {
	return strings.HasPrefix(href, "/") || strings.HasPrefix(href, baseHost)
}

func addPath(baseHost string, href string) string {
	// If href is already a full URL, return it
	if strings.HasPrefix(href, baseHost) {
		return href
	}

	// Remove leading slash from href and trailing slash from baseHost
	href = strings.TrimPrefix(href, "/")
	baseHost = strings.TrimSuffix(baseHost, "/")

	return baseHost + "/" + href
}
