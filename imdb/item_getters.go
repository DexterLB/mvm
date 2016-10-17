package imdb

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"

	"github.com/jbowtie/gokogiri/xml"
	"github.com/kennygrant/sanitize"
)

// ID returns the item's IMDB ID
func (s *Item) ID() int {
	return s.id
}

// Type tells whether the item a movie, series or episode
func (s *Item) Type() (ItemType, error) {
	if s.itemType != Unknown {
		return s.itemType, nil
	}

	mainPage, err := s.page("combined")
	if err != nil {
		return -1, err
	}

	_, err = s.episodeTitle()
	if err == nil {
		s.itemType = Episode
		return Episode, nil
	}

	eplist, err := mainPage.Search(`//h5[text()='Seasons:']`)
	if err == nil && len(eplist) > 0 {
		s.itemType = Series
		return Series, nil
	}

	s.itemType = Movie
	return Movie, nil
}

// Title returns the item's title
func (s *Item) Title() (string, error) {
	if s.title != nil {
		return *s.title, nil
	}

	episodeTitle, err := s.episodeTitle()
	if err == nil {
		s.title = &episodeTitle
		return episodeTitle, nil
	}

	h1, err := s.firstMatching("combined", `//h1`)
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

// Year returns the item's airing year
func (s *Item) Year() (int, error) {
	if s.year != nil {
		return *s.year, nil
	}

	mainPage, err := s.page("combined")
	if err != nil {
		return 0, err
	}

	itemType, err := s.Type()
	if err != nil {
		return 0, err
	}

	var yearText string

	if itemType == Episode {
		var yearSpans []xml.Node
		yearSpans, err = mainPage.Search(`//h1//span`)
		if err != nil {
			return 0, err
		}

		if len(yearSpans) == 0 {
			return 0, fmt.Errorf("can't find year element")
		}

		yearText = strings.Trim(yearSpans[0].LastChild().Content(), " ()")
	} else {
		var yearLinks []xml.Node
		yearLinks, err = mainPage.Search(`//a[contains(@href,'/year/')]`)
		if err != nil {
			return 0, err
		}

		if len(yearLinks) == 0 {
			return 0, fmt.Errorf("can't find year element")
		}

		yearText = yearLinks[0].Content()
	}

	var year int
	year, err = strconv.Atoi(yearText)
	if err != nil {
		return 0, fmt.Errorf("year is not a number: %s", err)
	}

	s.year = &year

	return year, nil
}

// OtherTitles returns the item's alternative titles
func (s *Item) OtherTitles() (map[string]string, error) {
	releaseInfoPage, err := s.page("releaseinfo")
	if err != nil {
		return nil, err
	}

	pairElements, err := releaseInfoPage.Search(`//*[@id='akas']/tr`)
	if err != nil {
		return nil, fmt.Errorf("can't find title elements: %s", err)
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

// ReleaseDate returns the item's release date.
// Only applicable for Movie and Episode.
func (s *Item) ReleaseDate() (time.Time, error) {
	itemType, err := s.Type()
	if err != nil {
		return time.Time{}, err
	}

	var dateText string

	if itemType == Episode {
		info, err := s.episodeInfo()
		if err != nil {
			return time.Time{}, err
		}

		dateText = info[0]
	} else {
		releaseDateElement, err := s.firstMatching(
			"combined",
			`//div[preceding-sibling::h5[contains(text(),'Release Date')]]`,
		)
		if err != nil {
			return time.Time{}, fmt.Errorf("unable to find release date element: %s", err)
		}

		matcher := regexp.MustCompile(`([0-9]{1,2} [A-Z][a-z]* [0-9]{4})`)
		groups := matcher.FindStringSubmatch(releaseDateElement.Content())
		if len(groups) < 2 {
			return time.Time{}, fmt.Errorf("unable to find release date")
		}

		dateText = groups[1]
	}

	return parseDate(dateText)
}

// Tagline returns the slogan. Probably only applicable for Movie.
func (s *Item) Tagline() (string, error) {
	taglineElement, err := s.firstMatching(
		"combined",
		`//div[preceding-sibling::h5[text()='Tagline:']]`,
	)
	if err != nil {
		return "", fmt.Errorf("unable to find tagline element: %s", err)
	}

	return strings.Trim(sanitize.HTML(taglineElement.Content()), " \n\t"), nil
}

// Duration returns the item's duration (rounded to minutes).
// Probably only applicable to Movie and Episode
func (s *Item) Duration() (time.Duration, error) {
	durationElement, err := s.firstMatching(
		"combined",
		`//div[preceding-sibling::h5[text()='Runtime:']]`,
	)
	if err != nil {
		return 0, fmt.Errorf("unable to find duration element: %s", err)
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

// Languages returns a slice with the names of languages for the item
func (s *Item) Languages() ([]*Language, error) {
	mainPage, err := s.page("combined")
	if err != nil {
		return nil, err
	}

	languageElements, err := mainPage.Search(
		`//div[preceding-sibling::h5[text()='Language:']]//a[contains(@href,'/language/')]`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to find language elements: %s", err)
	}

	matcher := regexp.MustCompile(`\/language\/(\w+)`)

	languages := make([]*Language, len(languageElements))
	for i := range languageElements {
		groups := matcher.FindStringSubmatch(
			languageElements[i].Attribute("href").String(),
		)
		if len(groups) < 2 {
			return nil, fmt.Errorf("invalid language")
		}
		lang, err := language.ParseBase(groups[1])
		if err != nil {
			return nil, fmt.Errorf("invalid language: %s", err)
		}
		languages[i] = (*Language)(&lang)
	}

	return languages, nil
}

// Plot returns the item's short plot summary
func (s *Item) Plot() (string, error) {
	plotElement, err := s.firstMatching(
		"combined",
		`//div[@class='info-content' and preceding-sibling::h5[text()='Plot:']]/text()`,
	)
	if err != nil {
		return "", fmt.Errorf("unable to find plot element: %s", err)
	}
	return strings.Trim(plotElement.Content(), " \t\n"), nil
}

// PlotMedium returns the item's medium-sized plot (summary)
func (s *Item) PlotMedium() (string, error) {
	summaryElement, err := s.firstMatching(
		"plotsummary",
		`//p[@class='plotSummary']`,
	)
	if err != nil {
		return "", fmt.Errorf("unable to find medium plot element: %s", err)
	}

	return strings.Trim(summaryElement.Content(), " \t\n"), nil
}

// PlotLong returns the item's long synopsis of the plot
func (s *Item) PlotLong() (string, error) {
	synopsisElement, err := s.firstMatching(
		"synopsis",
		`//div[@id='swiki.2.1']`,
	)
	if err != nil {
		return "", fmt.Errorf("unable to find long plot element: %s", err)
	}

	return strings.Trim(synopsisElement.Content(), " \t\n"), nil
}

// PosterURL returns.. the item's Poster URL (jpg image)
func (s *Item) PosterURL() (string, error) {
	posterElement, err := s.firstMatching(
		"combined",
		`//a[@name='poster']/img`,
	)
	if err != nil {
		return "", fmt.Errorf("unable to find poster url element: %s", err)
	}

	src := posterElement.Attribute("src")
	if src == nil {
		return "", fmt.Errorf("malformed poster image")
	}

	url := src.String()

	firstMatcher := regexp.MustCompile(`^(http:.+@@)`)
	secondMatcher := regexp.MustCompile(`^(http:.+?)\.[^\/]+$`)

	groups := firstMatcher.FindStringSubmatch(url)
	if len(groups) >= 2 {
		return groups[1] + ".jpg", nil
	}

	groups = secondMatcher.FindStringSubmatch(url)
	if len(groups) >= 2 {
		return groups[1] + ".jpg", nil
	}

	// return "", fmt.Errorf("can't parse poster url: '%s'", url)
	return url, nil
}

// Rating returns the item's rating
func (s *Item) Rating() (float32, error) {
	ratingElement, err := s.firstMatching(
		"combined",
		`//*[@class='starbar-meta']/b`,
	)
	var rating float64

	if err != nil {
		ratingElement, err = s.firstMatching(
			"",
			`//div[@class='ratingValue']`,
		)

		if err != nil {
			return 0, fmt.Errorf("unable to find rating element: %s", err)
		}
	}

	matcher := regexp.MustCompile(`(\d.\d)\/10`)
	groups := matcher.FindStringSubmatch(ratingElement.Content())

	if len(groups) < 2 {
		return 0, fmt.Errorf("can't find rating")
	}

	rating, err = strconv.ParseFloat(groups[1], 32)
	if err != nil {
		return 0, fmt.Errorf("unable to parse rating: %s", err)
	}

	return float32(rating), nil
}

// Votes returns the item's rating's number of votes
func (s *Item) Votes() (int, error) {
	votesElement, err := s.firstMatching(
		"combined",
		`//div[@id='tn15rating']//a[@class='tn15more']`,
	)
	if err != nil {
		votesElement, err = s.firstMatching(
			"",
			`//div[@class='imdbRating']//span[@class='small']`,
		)
		if err != nil {
			return 0, fmt.Errorf("unable to find votes element: %s", err)
		}
	}

	matcher := regexp.MustCompile(`([\d,]+)`)
	groups := matcher.FindStringSubmatch(votesElement.Content())

	if len(groups) < 2 {
		return 0, fmt.Errorf("can't find votes")
	}

	votes, err := strconv.Atoi(strings.Replace(groups[1], `,`, ``, 5))
	if err != nil {
		return 0, fmt.Errorf("unable to parse votes: %s", err)
	}

	return votes, nil
}

// SeasonEpisode returns an episode's season and episode numbers
func (s *Item) SeasonEpisode() (int, int, error) {
	itemType, err := s.Type()
	if err != nil {
		return 0, 0, fmt.Errorf("can't determine item type: %s", err)
	}

	if itemType != Episode {
		return 0, 0, nil
	}

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
func (s *Item) Series() (*Item, error) {
	seriesLinkElement, err := s.firstMatching(
		"combined",
		`//div[preceding-sibling::h5[contains(text(),'TV Series:')]]/a`,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to find series link element: %s", err)
	}

	href := seriesLinkElement.Attribute("href")
	if href == nil {
		return nil, fmt.Errorf("malformed series link")
	}

	id, err := idFromLink(href.String())
	if err != nil {
		return nil, err
	}

	return NewWithClient(id, s.client), nil
}

// episodeTitle returns the title of an episode
func (s *Item) episodeTitle() (string, error) {
	titleElement, err := s.firstMatching("combined", `//h1//span//em`)
	if err != nil {
		return "", fmt.Errorf("item not an episode?: %s", err)
	}

	title := titleElement.InnerHtml()

	if title == "" {
		return "", fmt.Errorf("empty title")
	}

	return title, nil
}

// episodeInfo returns a text block containing the episode's air date and number
func (s *Item) episodeInfo() ([]string, error) {
	infoElement, err := s.firstMatching(
		"combined",
		`//div[@class='info-content' and preceding-sibling::h5[contains(text(),'Original Air Date')]]`,
	)

	if err != nil {
		return nil, fmt.Errorf("item not an episode?: %s", err)
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

// Seasons returns a slice of all seasons in this item (applicable for Series).
// Please note that the indices in the slice might have nothing to do with
// the respective season numbers. For that, call Number() on each season.
func (s *Item) Seasons() ([]*Season, error) {
	mainPage, err := s.page("combined")
	if err != nil {
		return nil, err
	}

	seasonElements, err := mainPage.Search(
		`//div[preceding-sibling::h5[text()='Seasons:']]/a[contains(@href,'episodes?season')]`,
	)
	if err != nil {
		return nil, fmt.Errorf("can't find season elements")
	}

	seasons := make([]*Season, len(seasonElements))

	for i := range seasonElements {
		link := strings.Trim(seasonElements[i].Content(), " \t\n")
		url := fmt.Sprintf(
			"http://akas.imdb.com/title/tt%07d/episodes?season=%s",
			s.ID(), link,
		)
		season := NewSeasonWithClient(url, s.client)

		seasons[i] = season
	}
	return seasons, nil
}
