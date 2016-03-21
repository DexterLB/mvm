package library

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Show is a movie or an episode of a series
type Show struct {
	gorm.Model
	CommonData

	ReleaseDate time.Time `json:"release_date"`
	Tagline     string    `json:"tagline"`

	EpisodeData *EpisodeData `json:"episode_data"`

	Files []*VideoFile `json:"files",gorm:"ForeignKey:ShowID"`
}

// CommonData contains fields shared by movies, episodes and series
type CommonData struct {
	ImdbID      int               `json:"imdb_id",sql:"unique"`
	Title       string            `json:"title"`
	Year        int               `json:"year"`
	OtherTitles map[string]string `sql:"-",gorm:"-",json:"other_titles"`
	Duration    Duration          `json:"duration"`
	Plot        string            `json:"plot"`
	PlotMedium  string            `json:"plot_medium"`
	PlotLong    string            `json:"plot_long"`
	PosterURL   string            `json:"poster_url"`
	ImdbRating  float32           `json:"imdb_rating"`
	ImdbVotes   int               `json:"imdb_votes"`
	Languages   []string          `sql:"-",gorm:"-",json:"languages"`

	// No need to touch those - they are updated upon contact with the DB
	OtherTitlesJSON string `json:"other_titles_raw"`
	LanguagesJSON   string `json:"languages_raw"`
}

// EpisodeData contains episode-specific keys
type EpisodeData struct {
	ID       uint `gorm:"primary_key"`
	Season   int  `json:"season"`
	Episode  int  `json:"episode"`
	SeriesID uint
}

// Series represents a series
type Series struct {
	gorm.Model
	CommonData

	Episodes []*Show `json:"episodes",gorm:"ForeignKey:SeriesID"`
}

// VideoFile reprsesents a file for an episode or movie
type VideoFile struct {
	gorm.Model

	Filename   string   `json:"filename"`
	FileSize   uint     `json:"filesize"`
	Resolution [2]uint  `json:"resolution"`
	OsdbHash   uint     `json:"osdb_hash"`
	Format     string   `json:"format"`
	Duration   Duration `json:"duration"`

	LastPlayed   time.Time `json:"last_played"`
	LastPosition Duration  `json:"last_position"`

	ShowID uint
}

type Duration time.Duration

func (d *Duration) Scan(src interface{}) error {
	switch v := src.(type) {
	case int64:
		*d = Duration(Duration(v))
	default:
		return fmt.Errorf("unknown duration type")
	}
	return nil
}

func (c *CommonData) onSave() error {
	var err error
	c.LanguagesJSON, err = marshalString(&c.Languages)
	if err != nil {
		return fmt.Errorf("cannot convert languages to json: %s", err)
	}
	c.OtherTitlesJSON, err = marshalString(&c.OtherTitles)
	if err != nil {
		return fmt.Errorf("cannot convert other titles to json: %s", err)
	}
	return err
}

func (c *CommonData) onLoad() error {
	c.OtherTitles = make(map[string]string)
	if c.OtherTitlesJSON != "" {
		err := json.Unmarshal([]byte(c.OtherTitlesJSON), &c.OtherTitles)
		if err != nil {
			return fmt.Errorf("cannot parse other titles: %s", err)
		}
	}

	if c.LanguagesJSON != "" {
		err := json.Unmarshal([]byte(c.LanguagesJSON), &c.Languages)
		if err != nil {
			return fmt.Errorf("cannot parse languages: %s", err)
		}
	}
	return nil
}

func marshalString(value interface{}) (string, error) {
	data, err := json.Marshal(value)
	return string(data), err
}
