package library

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Show is a movie or an episode of a series
type Show struct {
	gorm.Model
	CommonData
	EpisodeData

	ReleaseDate time.Time `json:"release_date"`
	Tagline     string    `json:"tagline"`

	Files []*VideoFile `json:"files",gorm:"ForeignKey:ShowID"`
}

// CommonData contains fields shared by movies, episodes and series
type CommonData struct {
	ImdbID      int             `json:"imdb_id",sql:"unique"`
	Title       string          `json:"title"`
	Year        int             `json:"year"`
	OtherTitles MapStringString `gorm:"type:blob",json:"other_titles"`
	Duration    Duration        `gorm:"type:integer",json:"duration"`
	Plot        string          `json:"plot"`
	PlotMedium  string          `json:"plot_medium"`
	PlotLong    string          `json:"plot_long"`
	PosterURL   string          `json:"poster_url"`
	ImdbRating  float32         `json:"imdb_rating"`
	ImdbVotes   int             `json:"imdb_votes"`
	Languages   Languages       `gorm:"type:text",json:"languages"`
}

// EpisodeData contains episode-specific keys
type EpisodeData struct {
	Season   int `json:"season"`
	Episode  int `json:"episode"`
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

	Path        string   `json:"filename",sql:"unique"`
	Size        uint64   `json:"filesize"`
	ResolutionX uint     `json:"resolution"`
	ResolutionY uint     `json:"resolution"`
	OsdbHash    uint64   `json:"osdb_hash"`
	Format      string   `json:"format"`
	Duration    Duration `json:"duration"`

	LastPlayed   time.Time `json:"last_played"`
	LastPosition Duration  `json:"last_position"`

	ShowID uint
}
