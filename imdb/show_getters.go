package imdb

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
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

	_, err = s.episodeTitle()
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

	episodeTitle, err := s.episodeTitle()
	if err == nil {
		s.title = &episodeTitle
		return episodeTitle, nil
	}

	h1, err := firstMatching(s.mainPage, `//h1`)
	if err != nil {
		return "", fmt.Errorf("unable to find title element: %s", err)
	}

	title := strings.Split(h1.InnerHtml(), "<span")[0]
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
func (s *Show) OtherTitles() (map[string]string, error) {
	releaseInfoPage, err := s.releaseInfoPage()
	if err != nil {
		return nil, err
	}

	pairElements, err := releaseInfoPage.Search(`//*[@id='akas']/tr`)
	if err != nil {
		return nil, fmt.Errorf("can't find title elements")
	}

	titles := make(map[string]string)
	for i := range pairElements {
		versionElements, err := pairElements[i].Search(`td[position()=1]`)
		if err != nil {
			return nil, fmt.Errorf("can't parse title version element: %s", err)
		}
		titleElements, err := pairElements[i].Search(`td[position()=2]`)
		if err != nil {
			return nil, fmt.Errorf("can't parse title element: %s", err)
		}
		if len(versionElements) == 0 || len(titleElements) == 0 {
			return nil, fmt.Errorf("can't find title")
		}

		titles[versionElements[0].Content()] = titleElements[0].Content()
	}

	return titles, nil
}

// ReleaseDate returns the show's release date.
// Only applicable for Movie and Episode.
func (s *Show) ReleaseDate() (*time.Time, error) {
	showType, err := s.Type()
	if err != nil {
		return nil, err
	}

	var dateText string

	if showType == Episode {
		info, err := s.episodeInfo()
		if err != nil {
			return nil, err
		}

		dateText = info[0]
	} else {
		releaseDateElement, err := firstMatching(
			s.mainPage,
			`//div[preceding-sibling::h5[contains(text(),'Release Date')]]`,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to find release date element: %s", err)
		}

		matcher := regexp.MustCompile(`([0-9]{1,2} [A-Z][a-z]* [0-9]{4})`)
		groups := matcher.FindStringSubmatch(releaseDateElement.Content())
		if len(groups) < 2 {
			return nil, fmt.Errorf("unable to find release date")
		}

		dateText = groups[1]
	}

	return parseDate(dateText)
}

// Tagline returns the slogan. Probably only applicable for Movie.
func (s *Show) Tagline() (string, error) {
	taglineElement, err := firstMatching(
		s.mainPage,
		`//div[preceding-sibling::h5[text()='Tagline:']]`,
	)
	if err != nil {
		return "", err
	}

	return strings.Trim(sanitize.HTML(taglineElement.Content()), " \n\t"), nil
}

// Duration returns the show's duration (rounded to minutes).
// Probably only applicable to Movie and Episode
func (s *Show) Duration() (time.Duration, error) {
	durationElement, err := firstMatching(
		s.mainPage,
		`//div[preceding-sibling::h5[text()='Runtime:']]`,
	)
	if err != nil {
		return 0, err
	}

	matcher := regexp.MustCompile(`(\d+) min`)
	groups := matcher.FindStringSubmatch(durationElement.Content())
	if len(groups) < 2 {
		return 0, fmt.Errorf("unable to find duration")
	}

	minutes, err := strconv.Atoi(groups[1])
	if err != nil {
		return 0, fmt.Errorf("unable to parse duration: %s", err)
	}

	return time.Minute * time.Duration(minutes), nil
}

func (s *Show) Plot(level int) (string, error) {
	return "", fmt.Errorf("dummy method")
}

func (s *Show) PosterURL() (string, error) {
	return "", fmt.Errorf("dummy method")
}

// Rating returns the show's rating
func (s *Show) Rating() (float32, error) {
	ratingElement, err := firstMatching(
		s.mainPage,
		`//*[@class='starbar-meta']/b`,
	)
	if err != nil {
		return 0, err
	}

	matcher := regexp.MustCompile(`(\d.\d)\/10`)
	groups := matcher.FindStringSubmatch(ratingElement.Content())

	if len(groups) < 2 {
		return 0, fmt.Errorf("can't find rating")
	}

	log.Printf(groups[1])
	rating, err := strconv.ParseFloat(groups[1], 32)
	if err != nil {
		return 0, fmt.Errorf("unable to parse rating: %s", err)
	}

	return float32(rating), nil
}

func (s *Show) Votes() (int, error) {
	return 0, fmt.Errorf("dummy method")
}

// SeasonEpisode returns an episode's season and episode numbers
func (s *Show) SeasonEpisode() (int, int, error) {
	if s.season != nil && s.episode != nil {
		return *s.season, *s.episode, nil
	}

	info, err := s.episodeInfo()
	if err != nil {
		return 0, 0, err
	}

	matcher := regexp.MustCompile(`Season (\d+).*Episode (\d+)`)
	groups := matcher.FindStringSubmatch(info[len(info)-1])

	if len(groups) < 2 {
		return 0, 0, fmt.Errorf("can't find season/episode number")
	}

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

// Series returns the series this episode belongs to
func (s *Show) Series() (*Show, error) {
	seriesLinkElement, err := firstMatching(
		s.mainPage,
		`//div[preceding-sibling::h5[contains(text(),'TV Series:')]]/a`,
	)
	if err != nil {
		return nil, err
	}

	href := seriesLinkElement.Attribute("href")
	if href == nil {
		return nil, fmt.Errorf("malformed series link")
	}

	id, err := idFromLink(href.String())
	if err != nil {
		return nil, err
	}

	return New(id), nil
}

// episodeTitle returns the title of an episode
func (s *Show) episodeTitle() (string, error) {
	titleElement, err := firstMatching(s.mainPage, `//h1//span//em`)
	if err != nil {
		return "", fmt.Errorf("show not an episode?: %s", err)
	}

	title := titleElement.InnerHtml()

	if title == "" {
		return "", fmt.Errorf("empty title")
	}

	return title, nil
}

// episodeInfo returns a text block containing the episode's air date and number
func (s *Show) episodeInfo() ([]string, error) {
	infoElement, err := firstMatching(
		s.mainPage,
		`//div[@class='info-content' and preceding-sibling::h5[contains(text(),'Original Air Date')]]`,
	)

	if err != nil {
		return nil, fmt.Errorf("show not an episode?: %s", err)
	}

	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(infoElement.Content()))
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, nil
}
