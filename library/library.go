package library

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Library is a searchable library of movies, series and episodes
type Library struct {
	db *gorm.DB
}

// New creates a library connected to the specified database
func New(dbDriver string, arguments ...interface{}) (*Library, error) {
	db, err := gorm.Open(dbDriver, arguments...)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Show{}, &EpisodeData{}, &Series{}, &VideoFile{}, &Subtitle{})

	return &Library{
		db: db,
	}, nil
}

// HasSeriesWithImdbID checks if there exists a series with this id in the library
func (lib *Library) HasSeriesWithImdbID(id int) (bool, error) {
	err := lib.db.Where("imdb_id = ?", id).First(&Series{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetSeriesByImdbID finds the series by its imdb id, creating it if it doesn't exist
func (lib *Library) GetSeriesByImdbID(id int) (*Series, error) {
	series := &Series{}
	err := lib.db.Where("imdb_id = ?", id).FirstOrCreate(series).Error
	if err != nil {
		return nil, err
	}

	series.ImdbID = id
	err = lib.db.Model(series).Association("Episodes").Find(&series.Episodes).Error
	if err != nil {
		return nil, err
	}

	return series, err
}

// GetSeriesByEpisode returns the series this show belongs to (if it's an episode)
// or nil otherwise
func (lib *Library) GetSeriesByEpisode(episode *Show) (*Series, error) {
	series := &Series{}
	err := lib.db.Where("id = ?", episode.SeriesID).First(series).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return series, nil
}

// HasShowWithImdbID checks if there exists a show with this id in the library
func (lib *Library) HasShowWithImdbID(id int) (bool, error) {
	err := lib.db.Where("imdb_id = ?", id).First(&Show{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetShowByImdbID finds the show by its imdb id, creating it if it doesn't exist
func (lib *Library) GetShowByImdbID(id int) (*Show, error) {
	if id == 0 {
		return nil, fmt.Errorf("imdb id can't be 0")
	}

	show := &Show{}
	err := lib.db.Where("imdb_id = ?", id).FirstOrCreate(show).Error
	if err != nil {
		return nil, err
	}

	show.ImdbID = id

	err = lib.db.Model(show).Association("Files").Find(&show.Files).Error
	if err != nil {
		return nil, err
	}

	err = lib.db.Model(show).Association("Subtitles").Find(&show.Subtitles).Error
	if err != nil {
		return nil, err
	}

	return show, err
}

// HasFileWithPath checks if there exists a file with this path in the library
func (lib *Library) HasFileWithPath(path string) (bool, error) {
	err := lib.db.Where("path = ?", path).First(&VideoFile{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetFileByPath finds the file by its path, creating it if it doesn't exist
func (lib *Library) GetFileByPath(path string) (*VideoFile, error) {
	if path == "" {
		return nil, fmt.Errorf("path can't be empty")
	}

	file := &VideoFile{}
	err := lib.db.Where("path = ?", path).FirstOrCreate(file).Error
	if err != nil {
		return nil, err
	}
	file.Path = path
	return file, err
}

// GetSubtitleByHash finds the subtitle by its hash, creating it if it doesn't exist
func (lib *Library) GetSubtitleByHash(hash string) (*Subtitle, error) {
	subtitle := &Subtitle{}
	err := lib.db.Where("hash = ?", hash).FirstOrCreate(subtitle).Error
	if err != nil {
		return nil, err
	}
	subtitle.Hash = hash
	return subtitle, err
}

// Save saves the item to the library
func (lib *Library) Save(item interface{}) error {
	return lib.db.Save(item).Error
}
