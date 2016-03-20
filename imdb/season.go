package imdb

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	htmlParser "github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
)

// Season represents a single season from a series
type Season struct {
	url          string
	seasonNumber *int
	document     *htmlParser.HtmlDocument
	episodes     []*Item
}

// NewSeason creates a season from its url
func NewSeason(url string) *Season {
	return &Season{
		url: url,
	}
}

// URL returns the sason's url
func (s *Season) URL() string {
	return s.url
}

// Number returns the sason's number
func (s *Season) Number() (int, error) {
	if s.seasonNumber != nil {
		return *s.seasonNumber, nil
	}

	matcher := regexp.MustCompile(`episodes\?season=(\d+)`)
	groups := matcher.FindStringSubmatch(s.URL())

	if len(groups) < 2 {
		return 0, fmt.Errorf("invalid season url")
	}

	number, err := strconv.Atoi(groups[1])
	if err != nil {
		return 0, fmt.Errorf("unable to parse season number: %s", err)
	}

	s.seasonNumber = &number

	return number, nil
}

// Episodes returns an ordered slice of all episodes in this season.
// The returned items will have be of type Episode, and their Title and
// SeasonEpisode methods will return pre-cached results.
//
// Please note that although the episodes will probably be in the order
// they've come out, you shouldn't count on the fact that episode numbers
// have anything to do with indices in the slice. If you need the episode
// number, call SeasonEpisode on the episode itself.
func (s *Season) Episodes() ([]*Item, error) {
	if s.episodes != nil {
		return s.episodes, nil
	}

	page, err := s.page()
	if err != nil {
		return nil, err
	}

	episodeElements, err := page.Search(
		`//div[contains(@class,'eplist')]//div[contains(@itemprop,'episode')]`,
	)
	if err != nil {
		return nil, fmt.Errorf("can't find episode elements: %s", err)
	}

	episodes := make([]*Item, len(episodeElements))

	idMatcher := regexp.MustCompile(`tt([0-9]+)`)

	for i := range episodeElements {
		numberElement, err := firstMatchingOnNode(
			episodeElements[i],
			`.//meta[contains(@itemprop,'episodeNumber')]`,
		)
		if err != nil {
			return nil, fmt.Errorf("episode without number: %s", err)
		}
		number, err := strconv.Atoi(numberElement.Attribute("content").String())
		if err != nil {
			return nil, fmt.Errorf("unable to parse episode number: %s", err)
		}

		link, err := firstMatchingOnNode(
			episodeElements[i],
			`.//a[contains(@itemprop,'name')]`,
		)
		if err != nil {
			return nil, fmt.Errorf("episode without link: %s", err)
		}
		groups := idMatcher.FindStringSubmatch(link.Attribute("href").String())
		if len(groups) < 2 {
			return nil, fmt.Errorf("unable to find episode id")
		}
		id, err := strconv.Atoi(groups[1])
		if err != nil {
			return nil, fmt.Errorf("unable to parse episode id: %s", err)
		}

		title := strings.Trim(link.Content(), " \n\t")

		seasonNumber, err := s.Number()
		if err != nil {
			return nil, err
		}

		episodes[i] = &Item{
			id:       id,
			title:    &title,
			itemType: Episode,
			season:   &seasonNumber,
			episode:  &number,
		}
	}

	s.episodes = episodes

	return episodes, nil
}

func (s *Season) page() (*xml.ElementNode, error) {
	if s.document == nil {
		page, err := parsePage(s.URL())
		if err != nil {
			return nil, err
		}
		s.document = page
	}

	return s.document.Root(), nil
}
