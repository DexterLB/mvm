package imdb

import (
	"fmt"
	"regexp"
	"strconv"

	htmlParser "github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
)

// Season represents a single season from a series
type Season struct {
	url          string
	seasonNumber *int
	document     *htmlParser.HtmlDocument
	episodes     map[int]*Show
}

// NewSeason creates a season from its url
func NewSeason(url string) *Season {
	return &Season{
		url: url,
	}
}

// Url returns the sason's url
func (s *Season) Url() string {
	return s.url
}

// Number returns the sason's number
func (s *Season) Number() (int, error) {
	if s.seasonNumber != nil {
		return *s.seasonNumber, nil
	}

	matcher := regexp.MustCompile(`episodes\?season=(\d+)`)
	groups := matcher.FindStringSubmatch(s.Url())

	if len(groups) < 2 {
		return 0, fmt.Errorf("invalid season url")
	}

	number, err := strconv.Atoi(groups[1])
	if err != nil {
		return 0, fmt.Errorf("unable to parse season number: %s", err)
	}

	return number, nil
}

func (s *Season) page() (*xml.ElementNode, error) {
	if s.document == nil {
		page, err := parsePage(s.Url())
		if err != nil {
			return nil, err
		}
		s.document = page
	}

	return s.document.Root(), nil
}
