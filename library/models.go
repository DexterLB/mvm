package mvm

import (
	"time"

	"golang.org/x/text/language"

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
	Year        uint              `json:"year"`
	OtherTitles map[string]string `json:"other_titles"`
	Duration    time.Duration     `json:"duration"`
	Plot        string            `json:"plot"`
	PlotMedium  string            `json:"plot_medium"`
	PlotLong    string            `json:"plot_long"`
	PosterURL   string            `json:"poster_url"`
	ImdbRating  float32           `json:"imdb_rating"`
	ImdbVotes   int               `json:"imdb_votes"`
	Languages   []language.Base   `json:"languages"`
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

	Filename   string        `json:"filename"`
	FileSize   uint          `json:"filesize"`
	Resolution [2]uint       `json:"resolution"`
	OsdbHash   uint          `json:"osdb_hash"`
	Format     string        `json:"format"`
	Duration   time.Duration `json:"duration"`

	LastPlayed   time.Time     `json:"last_played"`
	LastPosition time.Duration `json:"last_position"`

	ShowID uint
}
