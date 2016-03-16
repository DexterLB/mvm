package imdb

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/moovweb/gokogiri"
	htmlParser "github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
)

// Show represents a single show (either a movie or an episode)
type Show struct {
	id               int
	title            *string
	showType         ShowType
	season           *int
	episode          *int
	mainPageDocument *htmlParser.HtmlDocument
}

//go:generate stringer -type=ShowType
// ShowType is one of Unknown, Movie, Series and Episode
type ShowType int

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

func (s *Show) parsePage(name string) (*htmlParser.HtmlDocument, error) {
	data, err := s.openPage(name)
	if err != nil {
		return nil, err
	}

	page, err := gokogiri.ParseHtml(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing html: %s", err)
	}
	return page, nil
}

func (s *Show) openPage(name string) ([]byte, error) {
	url := fmt.Sprintf("http://akas.imdb.com/title/tt%07d/%s", s.id, name)

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
