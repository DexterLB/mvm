package imdb

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/moovweb/gokogiri/xml"
	_ "github.com/orchestrate-io/dvr"
)

// ID returns the show's IMDB ID
func (s *Show) ID() int {
	return s.id
}

// Type tells whether the show a movie, series or episode
func (s *Show) Type() (ShowType, error) {
	if s.showType != Unknown {
		return s.showType, nil
	}

	mainPage, err := s.mainPage()
	if err != nil {
		return -1, err
	}

	_, err = episodeTitle(mainPage)
	if err == nil {
		s.showType = Episode
		return Episode, nil
	}

	eplist, err := mainPage.Search(`//h5[text()='Seasons:']`)
	if err == nil && len(eplist) > 0 {
		s.showType = Series
		return Series, nil
	}

	s.showType = Movie
	return Movie, nil
}

// Title returns the show's title
func (s *Show) Title() (string, error) {
	if s.title != nil {
		return *s.title, nil
	}

	mainPage, err := s.mainPage()
	if err != nil {
		return "", err
	}

	episodeTitle, err := episodeTitle(mainPage)
	if err == nil {
		s.title = &episodeTitle
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

	title = strings.Trim(title, " \n\r\t\"")
	s.title = &title
	return title, nil
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
	if s.season != nil && s.episode != nil {
		return *s.season, *s.episode, nil
	}

	mainPage, err := s.mainPage()
	if err != nil {
		return 0, 0, err
	}

	info, err := episodeInfo(mainPage)
	if err != nil {
		return 0, 0, err
	}

	matcher := regexp.MustCompile(`Season (\d+).*Episode (\d+)`)
	groups := matcher.FindStringSubmatch(info[len(info)-1])

	season, err := strconv.Atoi(groups[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid season number: %s", groups[1])
	}

	episode, err := strconv.Atoi(groups[2])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid episode number: %s", groups[2])
	}

	s.season = &season
	s.episode = &episode

	return season, episode, nil
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

// episodeTitle returns the title of an episode
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

// episodeInfo returns a text block containing the episode's air date and number
func episodeInfo(mainPage *xml.ElementNode) ([]string, error) {
	infoElements, err := mainPage.Search(
		`//div[@class='info-content' and preceding-sibling::h5[contains(text(),'Original Air Date')]]`,
	)
	if err != nil {
		return nil, err
	}

	if len(infoElements) == 0 {
		return nil, fmt.Errorf("unable to find info element (show not an episode?)")
	}

	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(infoElements[0].Content()))
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, nil
}
