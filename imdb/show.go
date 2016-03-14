package imdb

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/moovweb/gokogiri"
	htmlParser "github.com/moovweb/gokogiri/html"
	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
)

// Show represents a single show (either a movie or an episode)
type Show struct {
	id               int
	mainPageDocument *htmlParser.HtmlDocument
}

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

// Title returns the show's title
func (s *Show) Title() (string, error) {
	mainPage, err := s.mainPage()
	if err != nil {
		return "", err
	}

	h1, err := mainPage.Search(xpath.Compile("h1"))
	if err != nil {
		return "", err
	}

	if len(h1) == 0 {
		return "", fmt.Errorf("unable to find title element")
	}

	title := strings.Split(h1[0].InnerHtml(), "<span")[0]
	if title == "" {
		return "", fmt.Errorf("empty title")
	}

	return title, nil
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
