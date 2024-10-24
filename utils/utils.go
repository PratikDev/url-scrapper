package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func GetValidURLs(wg *sync.WaitGroup, mut *sync.Mutex, results *[]string, seen *[]string, url string) {
	defer wg.Done()

	// if url is already seen, return
	if isVisited(mut, seen, url) {
		return
	}

	statusCode, node := fetchURL(url)

	// if the status is not ok, return
	if statusCode != http.StatusOK {
		return
	}

	mut.Lock()
	*results = append(*results, url)
	mut.Unlock()

	hrefs := getHrefs(url, node)
	for _, href := range hrefs {
		wg.Add(1)
		// recursively call getValidURLs
		go GetValidURLs(wg, mut, results, seen, href)
	}
}

func isVisited(mut *sync.Mutex, seen *[]string, url string) bool {
	mut.Lock()
	defer mut.Unlock()
	if slices.Contains(*seen, url) {
		return true
	}
	*seen = append(*seen, url)
	return false
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

func fetchURL(url string) (statusCode int, node *html.Node) {
	fmt.Println("Checking URL:", url)

	response, err := http.Get(url)
	if err != nil {
		return http.StatusServiceUnavailable, nil
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
