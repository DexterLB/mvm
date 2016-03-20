// Package imdb provides easy access to publicly available data on IMDB.
// Items are accessed by their IMDB ID, and all getter methods called
// on them are lazy (an http request will be made only when data is needed,
// and this will happen only once). There is also a convenience AllData()
// method, which fetches all available data at once.
package imdb
