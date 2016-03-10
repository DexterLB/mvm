package mvm

type EpisodeData struct {
	SeriesTitle  string `json:"series_title"`
	SeriesYear   uint   `json:"series_year"`
	SeriesImdbID string `json:"series_imdb_id"`
	Season       int    `json:"season"`
	Episode      int    `json:"episode"`
}

type Movie struct {
	Episode *EpisodeData `json:"episode_data"`
	Title   string       `json:"title"`
	Year    uint         `json:"year"`
	ImdbID  string       `json:"imdb_id"`
}
