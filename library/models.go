package library

import (
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

// Show is a movie or an episode of a series
type Show struct {
	gorm.Model
	sync.Mutex
	CommonData
	EpisodeData

	ReleaseDate time.Time `json:"release_date"`
	Tagline     string    `json:"tagline"`

	Files     []*VideoFile `json:"files",gorm:"ForeignKey:ShowID"`
	Subtitles []*Subtitle  `json:"subtitles",gorm:"ForeignKey:ShowID"`
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

	ImdbError *string `json:"imdb_error"`
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
	sync.Mutex
	CommonData

	Episodes []*Show `json:"episodes",gorm:"ForeignKey:SeriesID"`
}

// VideoFile reprsesents a file for an episode or movie
type VideoFile struct {
	gorm.Model
	sync.Mutex

	Path         string    `json:"filename",sql:"unique"`
	Size         uint64    `json:"filesize"`
	ResolutionX  uint      `json:"resolution"`
	ResolutionY  uint      `json:"resolution"`
	OsdbHash     BigUint64 `gorm:"type:varchar(16)",json:"osdb_hash"`
	VideoFormat  string    `json:"video_format"`
	AudioFormat  string    `json:"audio_format"`
	Framerate    float32   `json:"framerate"`
	VideoBitrate float32   `json:"video_bitrate"`
	AudioBitrate float32   `json:"audio_bitrate"`
	Duration     Duration  `json:"duration"`

	LastPlayed   time.Time `json:"last_played"`
	LastPosition Duration  `json:"last_position"`

	ShowID uint

	ImportError *string `json:"import_error"`
	OsdbError   *string `json:"osdb_error"`
}

// Subtitle represents a subtitle file
type Subtitle struct {
	gorm.Model
	sync.Mutex

	Hash            string   `json:"hash"`
	Language        Language `gorm:"type:varchar(3)",json:"language"`
	HearingImpaired bool     `json:"hearing_impaired"`
	Filename        string   `json:"filename"`

	ShowID uint
}

// AfterCreate initializes values on an empty series
func (s *Series) AfterCreate() error {
	// nothing! FIXME: remove this function
	return nil
}

// AfterCreate initializes values on an empty show
func (s *Show) AfterCreate() error {
	// nothing! FIXME: remove this function
	return nil
}

// AfterCreate initializes values on an empty video file
func (v *VideoFile) AfterCreate() error {
	// nothing! FIXME: remove this function
	return nil
}
