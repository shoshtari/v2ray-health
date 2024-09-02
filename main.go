package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/shoshtari/v2ray-health/pkg"
)

func getPage(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get page")
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, pkg.ErrNotFound
		}
		return nil, fmt.Errorf("status was %d instead of 200", res.StatusCode)
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read data")
	}
	urls := strings.Split(string(resData), "\n")
	return urls, nil
}

func getUrls() ([]string, error) {
	subUrlTemplate := "https://raw.githubusercontent.com/barry-far/V2ray-Configs/main/Sub%d.txt"
	var urls []string

	subNumber := 1
	for true {
		subUrl := fmt.Sprintf(subUrlTemplate, subNumber+1)
		subNumber++
		res, err := getPage(subUrl)
		if errors.Is(err, pkg.ErrNotFound) {
			return urls, nil
		}
		if err != nil {
			panic(err)
		}
		urls = append(urls, res...)
	}
	return urls, nil

}

func removeDupUrls(urls []string) []string {
	var ans []string
	urlExist := make(map[string]struct{})
	for _, u := range urls {
		if _, exists := urlExist[u]; exists {
			continue
		}
		urlExist[u] = struct{}{}
		ans = append(ans, u)
	}
	return ans
}
func main() {
	urls, err := getUrls()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got %d urls\n", len(urls))
	urls = removeDupUrls(urls)
	fmt.Printf("Got %d unique urls\n", len(urls))
}
