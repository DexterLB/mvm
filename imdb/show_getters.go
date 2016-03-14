package imdb

import (
	"fmt"
	"strings"
	"time"

	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
)

// IsEpisode tells whether the show is part of a series
func (s *Show) IsEpisode() (bool, error) {
	mainPage, err := s.mainPage()
	if err != nil {
		return false, err
	}

	_, err = episodeTitle(mainPage)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// Title returns the show's title
func (s *Show) Title() (string, error) {
	mainPage, err := s.mainPage()
	if err != nil {
		return "", err
	}

	episodeTitle, err := episodeTitle(mainPage)
	if err == nil {
		return episodeTitle, nil
	}

	h1, err := mainPage.Search(xpath.Compile(`//h1`))
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

	return strings.Trim(title, " \n\r\t\""), nil
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

func episodeTitle(mainPage *xml.ElementNode) (string, error) {
	titleElements, err := mainPage.Search(xpath.Compile(`//h1//span//em`))
	if err != nil {
		return "", err
	}

	if len(titleElements) == 0 {
		return "", fmt.Errorf("unable to find title element (show not an episode?)")
	}

	return titleElements[0].InnerHtml(), nil
}
