package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const URL = "https://scrape-me.dreamsofcode.io"

func main() {
	response, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	rawHtml := string(bytes)
	node, err := html.Parse(strings.NewReader(rawHtml))
	if err != nil {
		panic(err)
	}

	hrefs := getHrefs(URL, node)
	fmt.Println(hrefs)
}

func getHrefs(baseHost string, n *html.Node) []string {
	result := []string{}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" && isBaseHost(baseHost, a.Val) {
				result = append(result, addPath(baseHost, a.Val))
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
	// Remove leading slash from href and trailing slash from baseHost
	href = strings.TrimPrefix(href, "/")
	baseHost = strings.TrimSuffix(baseHost, "/")
	return baseHost + "/" + href
}
