package jsonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DexterLB/mvm/imdb"
)

// Client connects to the IMDB jsonapi server and performs basic requests
type Client struct {
	HttpClient imdb.HttpGetPoster
	Address    string
}

// Item returns information for an item with a given IMDB ID
func (c *Client) Item(id int) (*imdb.ItemData, error) {
	resp, err := c.HttpClient.Get(fmt.Sprintf(
		"%s/item?id=%d", c.Address, id,
	))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: %s", resp.Status)
	}

	data := &imdb.ItemData{}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// SearchQuery performs a search on IMDB
func (c *Client) Search(query *imdb.SearchQuery) ([]*imdb.ShortItem, error) {
	request := &bytes.Buffer{}
	err := json.NewEncoder(request).Encode(query)
	if err != nil {
		return nil, fmt.Errorf("unable to encode query: %s", err)
	}

	resp, err := c.HttpClient.Post(
		fmt.Sprintf("%s/search", c.Address),
		"application/json",
		request,
	)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: %s", resp.Status)
	}

	var data []*imdb.ShortItem
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}
