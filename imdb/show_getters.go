package imdb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/moovweb/gokogiri/xml"
	_ "github.com/orchestrate-io/dvr"
)

//go:generate stringer -type=ShowType
type ShowType int

const (
	// Movie is the type of a show which is a movie
	Movie ShowType = iota
	// Movie is the type of a show which is a series
	Series
	// Episode is the type of a show which is an episode
	Episode
)

// Type tells whether the show a movie, series or episode
func (s *Show) Type() (ShowType, error) {
	mainPage, err := s.mainPage()
	if err != nil {
		return -1, err
	}

	_, err = episodeTitle(mainPage)
	if err == nil {
		return Episode, nil
	}

	eplist, err := mainPage.Search(`//h5[text()='Seasons:']`)
	if err == nil && len(eplist) > 0 {
		return Series, nil
	}

	return Movie, nil
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

	h1, err := mainPage.Search(`//h1`)
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
	mainPage, err := s.mainPage()
	if err != nil {
		return 0, err
	}

	showType, err := s.Type()
	if err != nil {
		return 0, err
	}

	var yearText string

	if showType == Episode {
		yearSpans, err := mainPage.Search(`//h1//span`)
		if err != nil {
			return 0, err
		}

		if len(yearSpans) == 0 {
			return 0, fmt.Errorf("can't find year element")
		}

		yearText = strings.Trim(yearSpans[0].LastChild().Content(), " ()")
	} else {
		yearLinks, err := mainPage.Search(`//a[contains(@href,'/year/')]`)
		if err != nil {
			return 0, err
		}

		if len(yearLinks) == 0 {
			return 0, fmt.Errorf("can't find year element")
		}

		yearText = yearLinks[0].Content()
	}

	year, err := strconv.Atoi(yearText)
	if err != nil {
		return 0, fmt.Errorf("year is not a number: %s", err)
	}

	return year, nil
}

// OtherTitles returns the show's alternative titles
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
	titleElements, err := mainPage.Search(`//h1//span//em`)
	if err != nil {
		return "", err
	}

	if len(titleElements) == 0 {
		return "", fmt.Errorf("unable to find title element (show not an episode?)")
	}

	title := titleElements[0].InnerHtml()

	if title == "" {
		return "", fmt.Errorf("empty title")
	}

	return title, nil
}
