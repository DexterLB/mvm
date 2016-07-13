package imdb

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/jbowtie/gokogiri/xml"
)

// SearchQuery represents a query for items (shows, series or episodes)
type SearchQuery struct {
	Query    string
	Year     int
	Category ItemType
	Exact    bool
}

// paremeters returns the IMDB URL parameters for the search
func (s *SearchQuery) parameters() *url.Values {
	parameters := &url.Values{}
	if s.Year != 0 {
		parameters.Set("q", fmt.Sprintf("%s (%d)", s.Query, s.Year))
	} else {
		parameters.Set("q", s.Query)
	}

	if s.Exact {
		parameters.Set("exact", "true")
	}

	if s.Category != Unknown {
		parameters.Set("s", "tt")
		switch s.Category {
		case Movie:
			parameters.Set("ttype", "ft")
		case Series:
			parameters.Set("ttype", "tv")
		case Episode:
			parameters.Set("ttype", "ep")
		}
	}

	return parameters
}

// encode returns a GET request parameter string
func (s *SearchQuery) encode() string {
	return s.parameters().Encode()
}

// Search executes the query
func Search(query *SearchQuery) ([]*Item, error) {
	return SearchWithClient(query, nil)
}

// SearchWithClient executes the query using the given HTTP client to communicate
// with IMDB.
func SearchWithClient(query *SearchQuery, client HttpGetter) ([]*Item, error) {
	searchPage, err := parsePage(client, fmt.Sprintf(
		"http://akas.imdb.com/find?%s", query.encode(),
	))
	if err != nil {
		return nil, fmt.Errorf("unable to parse search page: %s", err)
	}

	rows, err := searchPage.Search(
		`//table[contains(@class, 'findList')]//tr[contains(@class, 'findResult')]`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to find search results: %s", err)
	}

	var items []*Item

	for i := range rows {
		item, err := parseSearchResult(rows[i])
		if err != nil {
			return nil, fmt.Errorf("unable to parse search result: %s", err)
		}
		if item != nil {
			items = append(items, item)
		}
	}

	return items, nil
}

func parseSearchResult(row xml.Node) (*Item, error) {
	textElements, err := row.Search("td[@class='result_text']")
	if err != nil {
		return nil, fmt.Errorf("unable to find result text: %s", err)
	}
	if len(textElements) < 1 {
		return nil, fmt.Errorf("no result text")
	}

	id, err := getSearchResultID(textElements[0])
	if err != nil {
		return nil, err
	}

	if id == 0 {
		return nil, nil
	}

	if len(textElements[0].Content()) == 0 {
		return nil, nil
	}

	groups := regexp.MustCompile(
		`(.+)[\s]+\(([0-9]+)\)[\s]*`,
	).FindStringSubmatch(
		textElements[0].Content(),
	)

	if len(groups) < 3 {
		return nil, fmt.Errorf(
			"unable to parse result line (%s)", textElements[0].Content(),
		)
	}

	title := strings.Trim(groups[1], " \r\t\n")

	year, err := strconv.Atoi(groups[2])
	if err != nil {
		return nil, fmt.Errorf("unable to parse year (%s): %s", groups[3], err)
	}

	item := New(id)
	item.title = &title
	item.year = &year

	return item, nil
}

func getSearchResultID(resultTextElement xml.Node) (int, error) {
	linkElements, err := resultTextElement.Search("a")
	if err != nil {
		return 0, fmt.Errorf("unable to find result link: %s", err)
	}
	if len(linkElements) < 1 {
		return 0, fmt.Errorf("no result link")
	}

	href := linkElements[0].Attribute("href")
	if href == nil {
		return 0, fmt.Errorf("malformed result link")
	}

	if !strings.Contains(href.String(), "/title/tt") {
		return 0, nil // not a valid search result
	}

	return idFromLink(href.String())
}
