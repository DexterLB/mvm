package imdb

import (
	"fmt"
	"strings"
	"time"

	"github.com/moovweb/gokogiri/xpath"
)

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

// Year returns the show's airing year
func (s *Show) Year() (int, error) {
	return 0, fmt.Errorf("dummy method")
}

func (s *Show) OtherTitles() ([]string, error) {
	return nil, fmt.Errorf("dummy method")
}

func (s *Show) ReleaseDate() (*time.Time, error) {
	return nil, fmt.Errorf("dummy method")
}

func (s *Show) Tagline() (string, error) {
	return "", fmt.Errorf("dummy method")
}

func (s *Show) Duration() (*time.Duration, error) {
	return nil, fmt.Errorf("dummy method")
}

func (s *Show) Plot(level int) (string, error) {
	return "", fmt.Errorf("dummy method")
}

func (s *Show) PosterURL() (string, error) {
	return "", fmt.Errorf("dummy method")
}

func (s *Show) Rating() (float32, error) {
	return 0, fmt.Errorf("dummy method")
}

func (s *Show) Votes() (int, error) {
	return 0, fmt.Errorf("dummy method")
}

func (s *Show) IsEpisode() (bool, error) {
	return false, fmt.Errorf("dummy method")
}

func (s *Show) SeasonEpisode() (int, int, error) {
	return 0, 0, fmt.Errorf("dummy method")
}

func (s *Show) SeriesID() (int, error) {
	return 0, fmt.Errorf("dummy method")
}

func (s *Show) SeriesTitle() (string, error) {
	return "", fmt.Errorf("dummy method")
}

func (s *Show) SeriesYear() (int, error) {
	return 0, fmt.Errorf("dummy method")
}
