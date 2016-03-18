package imdb

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/moovweb/gokogiri"
	htmlParser "github.com/moovweb/gokogiri/html"
)

func parsePage(url string) (*htmlParser.HtmlDocument, error) {
	data, err := openPage(url)
	if err != nil {
		return nil, err
	}

	page, err := gokogiri.ParseHtml(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing html: %s", err)
	}
	return page, nil
}

func openPage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to download imdb page: %s", err)
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read imdb page: %s", err)
	}

	return data, nil
}
