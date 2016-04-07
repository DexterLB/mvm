package imdb

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"golang.org/x/text/language"
)

// ItemData holds all fields for a item
type ItemData struct {
	ID   int
	Type ItemType

	Title       string
	Year        int
	OtherTitles map[string]string
	Duration    time.Duration
	Plot        string
	PlotMedium  string
	PlotLong    string
	PosterURL   string
	Rating      float32
	Votes       int
	Languages   []language.Base

	// Movie and Episode-only fields
	ReleaseDate time.Time

	// Movie-only fields
	Tagline string

	// Episode-only fields
	SeasonNumber  int
	EpisodeNumber int
	Series        *Item

	// Series-only fields
	Seasons []*Season
}

// String returns the data in a human-readable form
func (s *ItemData) String() string {
	languageNames := make([]string, len(s.Languages))

	for i, language := range s.Languages {
		languageNames[i] = language.String()
	}
	text := fmt.Sprintf(`id: %d
type: %s
title: %s
year: %d
other titles:
%s
duration: %s
short plot: %s
medium plot: %s
long plot: %s
poster url: %s
rating: %.2g
votes: %dk
languages: %s`,
		s.ID,
		s.Type,
		s.Title,
		s.Year,
		humanReadableMap(s.OtherTitles),
		s.Duration,
		s.Plot,
		shorten(s.PlotMedium),
		shorten(s.PlotLong),
		s.PosterURL,
		s.Rating,
		s.Votes/1000,
		strings.Join(languageNames, ", "),
	)

	if s.Type == Movie || s.Type == Episode {
		text = fmt.Sprintf(`%s
release date: %s`,
			text,
			s.ReleaseDate.Format("2006-01-02"),
		)
	}

	switch s.Type {
	case Movie:
		text = fmt.Sprintf(`%s
tagline: %s`,
			text,
			s.Tagline,
		)
	case Episode:
		text = fmt.Sprintf(`%s
season number: %d
episode number: %d
series id: %07d`,
			text,
			s.SeasonNumber,
			s.EpisodeNumber,
			s.Series.ID(),
		)
	case Series:
		seasonNumbers := make([]string, len(s.Seasons))
		for i, season := range s.Seasons {
			number, _ := season.Number()
			seasonNumbers[i] = fmt.Sprintf("%d", number)
		}
		text = fmt.Sprintf(`%s
seasons: %s`,
			text,
			strings.Join(seasonNumbers, ", "),
		)
	}
	return text
}

func shorten(text string) string {
	sentences := strings.Split(text, ". ")
	if len(sentences) < 1 {
		return ""
	}
	return sentences[0] + "..."
}

func humanReadableMap(m map[string]string) string {
	lines := make([]string, len(m))
	var i int
	for key := range m {
		lines[i] = key
		i++
	}
	sort.Strings(lines)

	for i, key := range lines {
		lines[i] = fmt.Sprintf(" > %s -> %s", key, m[key])
	}
	return strings.Join(lines, "\n")
}

func (s *Item) fillPlot(data *ItemData) error {
	var err error
	data.Plot, err = s.Plot()
	if err != nil {
		return err
	}

	data.PlotMedium, err = s.PlotMedium()
	if err != nil {
		return err
	}

	data.PlotLong, err = s.PlotLong()
	if err != nil {
		return err
	}

	return nil
}

func (s *Item) fillSpecificData(data *ItemData) error {
	var err error

	if data.Type == Movie || data.Type == Episode {
		data.ReleaseDate, err = s.ReleaseDate()
		if err != nil {
			return err
		}
	}

	switch data.Type {
	case Movie:
		data.Tagline, _ = s.Tagline() // tagline is optional
	case Episode:
		data.SeasonNumber, data.EpisodeNumber, err = s.SeasonEpisode()
		if err != nil {
			return err
		}

		data.Series, err = s.Series()
		if err != nil {
			return err
		}
	case Series:
		data.Seasons, err = s.Seasons()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Item) fillCommonData(data *ItemData) error {
	var err error

	data.Type, err = s.Type()
	if err != nil {
		return err
	}

	data.Title, err = s.Title()
	if err != nil {
		return err
	}

	data.Year, err = s.Year()
	if err != nil {
		return err
	}

	data.OtherTitles, err = s.OtherTitles()
	if err != nil {
		return err
	}

	data.Duration, err = s.Duration()
	if err != nil {
		return err
	}

	data.PosterURL, err = s.PosterURL()
	if err != nil {
		return err
	}

	data.Rating, err = s.Rating()
	if err != nil {
		return err
	}

	data.Votes, err = s.Votes()
	if err != nil {
		return err
	}

	data.Languages, err = s.Languages()
	if err != nil {
		return err
	}

	return nil
}

// AllData fetches all possible fields and returns them
func (s *Item) AllData() (*ItemData, error) {
	s.PreloadAll()

	data := &ItemData{
		ID: s.ID(),
	}

	err := s.fillCommonData(data)
	if err != nil {
		return nil, err
	}

	err = s.fillPlot(data)
	if err != nil {
		return nil, err
	}

	err = s.fillSpecificData(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
