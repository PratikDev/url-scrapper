package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func GetHrefs(baseHost string, n *html.Node) []string {
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
		result = append(result, GetHrefs(baseHost, c)...)
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

func FetchURL(url string) (statusCode int, node *html.Node) {
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
