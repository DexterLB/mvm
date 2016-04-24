package imdb

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	htmlParser "github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
)

// Item represents a single item (either a movie or an episode)
type Item struct {
	id              int
	title           *string
	itemType        ItemType
	season          *int
	episode         *int
	cachedDocuments map[string]*htmlParser.HtmlDocument
	cacheLock       sync.Mutex
}

// ItemType is one of Unknown, Movie, Series and Episode
type ItemType int

//go:generate stringer -type=ItemType

const (
	// Unknown is a null item type
	Unknown ItemType = iota
	// Movie is the type of a item which is a movie
	Movie
	// Series is the type of a item which is a series
	Series
	// Episode is the type of a item which is an episode
	Episode
)

// New creates a item from an IMDB ID
func New(id int) *Item {
	return &Item{
		id:              id,
		cachedDocuments: make(map[string]*htmlParser.HtmlDocument),
	}
}

// Free frees all resources used by the parser. You must always call it
// after you finish reading the attributes
func (s *Item) Free() {
	for name := range s.cachedDocuments {
		s.cachedDocuments[name].Free()
		delete(s.cachedDocuments, name)
	}
}

// PreloadAll loads all pages needed for this item by making parallel
// requests to IMDB. All subsequent calls to methods will be fast
// (won't generate a http request)
func (s *Item) PreloadAll() {
	wg := sync.WaitGroup{}
	wg.Add(4)

	load := func(name string) {
		_, _ = s.page(name)
		wg.Done()
	}

	go load("combined")
	go load("releaseinfo")
	go load("plotsummary")
	go load("synopsis")

	wg.Wait()
}

// page returns the html contents of the page at
// http://akas.imdb.com/title/tt<s.ID>/<name>
func (s *Item) page(name string) (*xml.ElementNode, error) {
	document, ok := s.cachedDocuments[name]
	if !ok {
		var err error

		document, err = s.parsePage(name)
		if err != nil {
			return nil, err
		}
		s.cacheLock.Lock()
		s.cachedDocuments[name] = document
		s.cacheLock.Unlock()
	}

	return document.Root(), nil
}

func (s *Item) parsePage(name string) (*htmlParser.HtmlDocument, error) {
	url := fmt.Sprintf("http://akas.imdb.com/title/tt%07d/%s", s.id, name)
	return parsePage(url)
}

// idFromLink extracts an IMDB ID from a link
func idFromLink(link string) (int, error) {
	matcher := regexp.MustCompile(`\/tt([0-9]+)`)
	groups := matcher.FindStringSubmatch(link)

	if len(groups) <= 1 || groups[1] == "" {
		return 0, fmt.Errorf("invalid link: %s", link)
	}

	id, err := strconv.Atoi(groups[1])
	if err != nil {
		return 0, fmt.Errorf("invalid imdb id: %s", err)
	}

	return id, nil
}

// parseDate parses a date from IMDB's default format
func parseDate(text string) (time.Time, error) {
	t, err := time.Parse("2 January 2006", text)
	if err != nil {
		return time.Time{}, fmt.Errorf("can't parse date string '%s': %s", text, err)
	}
	return t, nil
}

// firstMatching obtains a root node by calling page(),
// and then finds its first child node which matches the xpath
func (s *Item) firstMatching(pageName string, xpath string) (xml.Node, error) {
	page, err := s.page(pageName)
	if err != nil {
		return nil, err
	}

	return firstMatchingOnNode(page, xpath)
}

// firstMatchingOnNode finds its first child node which matches the xpath
func firstMatchingOnNode(node xml.Node, xpath string) (xml.Node, error) {
	elements, err := node.Search(xpath)
	if err != nil {
		return nil, err
	}

	if len(elements) == 0 {
		return nil, fmt.Errorf("unable to find element")
	}

	return elements[0], nil
}
