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

// Show represents a single show (either a movie or an episode)
type Show struct {
	id                   int
	title                *string
	showType             ShowType
	season               *int
	episode              *int
	mainPageDocument     *htmlParser.HtmlDocument
	plotSummaryDocument  *htmlParser.HtmlDocument
	plotSynopsisDocument *htmlParser.HtmlDocument
	releaseInfoDocument  *htmlParser.HtmlDocument
}

// ShowType is one of Unknown, Movie, Series and Episode
type ShowType int

//go:generate stringer -type=ShowType

const (
	// Unknown is a null show type
	Unknown ShowType = iota
	// Movie is the type of a show which is a movie
	Movie
	// Series is the type of a show which is a series
	Series
	// Episode is the type of a show which is an episode
	Episode
)

// New creates a show from an IMDB ID
func New(id int) *Show {
	return &Show{
		id: id,
	}
}

// Free frees all resources used by the parser. You must always call it
// after you finish reading the attributes
func (s *Show) Free() {
	if s.mainPageDocument != nil {
		s.mainPageDocument.Free()
		s.mainPageDocument = nil
	}

	if s.releaseInfoDocument != nil {
		s.releaseInfoDocument.Free()
		s.releaseInfoDocument = nil
	}

	if s.plotSummaryDocument != nil {
		s.plotSummaryDocument.Free()
		s.plotSummaryDocument = nil
	}

	if s.plotSynopsisDocument != nil {
		s.plotSynopsisDocument.Free()
		s.plotSynopsisDocument = nil
	}
}

// PreloadAll loads all pages needed for this show by making parallel
// requests to IMDB. All subsequent calls to methods will be fast
// (won't generate a http request)
func (s *Show) PreloadAll() {
	wg := sync.WaitGroup{}
	wg.Add(4)

	load := func(getPage func() (*xml.ElementNode, error)) {
		_, _ = getPage()
		wg.Done()
	}

	go load(s.mainPage)
	go load(s.releaseInfoPage)
	go load(s.plotSummaryPage)
	go load(s.plotSynopsisPage)

	wg.Wait()
}

func (s *Show) mainPage() (*xml.ElementNode, error) {
	if s.mainPageDocument == nil {
		page, err := s.parsePage("combined")
		if err != nil {
			return nil, err
		}
		s.mainPageDocument = page
	}

	return s.mainPageDocument.Root(), nil
}

func (s *Show) releaseInfoPage() (*xml.ElementNode, error) {
	if s.releaseInfoDocument == nil {
		page, err := s.parsePage("releaseinfo")
		if err != nil {
			return nil, err
		}
		s.releaseInfoDocument = page
	}

	return s.releaseInfoDocument.Root(), nil
}

func (s *Show) plotSummaryPage() (*xml.ElementNode, error) {
	if s.plotSummaryDocument == nil {
		page, err := s.parsePage("plotsummary")
		if err != nil {
			return nil, err
		}
		s.plotSummaryDocument = page
	}

	return s.plotSummaryDocument.Root(), nil
}

func (s *Show) plotSynopsisPage() (*xml.ElementNode, error) {
	if s.plotSynopsisDocument == nil {
		page, err := s.parsePage("synopsis")
		if err != nil {
			return nil, err
		}
		s.plotSynopsisDocument = page
	}

	return s.plotSynopsisDocument.Root(), nil
}

func (s *Show) parsePage(name string) (*htmlParser.HtmlDocument, error) {
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

// firstMatching obtains a root node by calling pageGetter,
// and then finds its first child node which matches the xpath
func firstMatching(pageGetter func() (*xml.ElementNode, error), xpath string) (xml.Node, error) {
	page, err := pageGetter()
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
