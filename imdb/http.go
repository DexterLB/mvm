package imdb

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jbowtie/gokogiri"
	htmlParser "github.com/jbowtie/gokogiri/html"
)

type HttpGetter interface {
	Get(url string) (resp *http.Response, err error)
}

func parsePage(client HttpGetter, url string) (*htmlParser.HtmlDocument, error) {
	data, err := openPage(client, url)
	if err != nil {
		return nil, err
	}

	page, err := gokogiri.ParseHtml(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing html: %s", err)
	}
	return page, nil
}

func openPage(client HttpGetter, url string) ([]byte, error) {
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to download imdb page: %s", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read imdb page: %s", err)
	}

	return data, nil
}
