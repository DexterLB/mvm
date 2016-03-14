package imdb

import (
	"fmt"
	"strings"
)

func ExampleShow_Title() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	title, err := movie.Title()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("title: %s\n", title)

	// Output:
	// title: Nochnoy dozor
}

func ExampleShow_Title_episode() {
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()

	title, err := episode.Title()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("title: %s\n", title)

	// Output:
	// title: The Beast Below
}

func ExampleShow_Year() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	year, err := movie.Year()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("year: %d\n", year)

	// Output:
	// year: 2004
}

func ExampleShow_OtherTitles() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	otherTitles, err := movie.OtherTitles()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("other titles: %s\n", strings.Join(otherTitles, ", "))

	// Output:
	// other titles: Night Watch
}

func ExampleShow_ReleaseDate() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	releaseDate, err := movie.ReleaseDate()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("release date: %s\n", releaseDate)

	// Output:
	// release date: 2006-03-03
}

func ExampleShow_Tagline() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	tagline, err := movie.Tagline()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("tagline: %s\n", tagline)

	// Output:
	// tagline: All That Stands Between Light And Darkness Is The Night Watch.
}

func ExampleShow_Duration() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	duration, err := movie.Duration()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("duration: %s\n", duration)

	// Output:
	// duration: 1h 54m
}

func ExampleShow_Plot() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	shortPlot, err := movie.Plot(0)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("short plot: %s\n", shortPlot)

	// Output:
	// short plot: A fantasy-thriller set in present-day Moscow where the respective forces that control daytime and nighttime do battle.
}

func ExampleShow_PosterURL() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	posterURL, err := movie.PosterURL()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("poster url: %s\n", posterURL)

	// Output:
	// poster url: http://ia.media-imdb.com/images/M/MV5BMjE0Nzk0NDkyOV5BMl5BanBnXkFtZTcwMjkzOTkyMQ@@._V1_SX640_SY720_.jpg
}

func ExampleShow_Rating() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	rating, err := movie.Rating()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("rating: %.1g\n", rating)

	// Output:
	// rating: 6.5
}

func ExampleShow_Votes() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	votes, err := movie.Votes()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("votes: %d\n", votes)

	// Output:
	// votes: 46450
}

func ExampleShow_IsEpisode() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()

	isMovieEpisode, err := movie.IsEpisode()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	isEpisodeEpisode, err := episode.IsEpisode()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("Is Nochnoy Dozor an episode? %v\n", isMovieEpisode)
	fmt.Printf("Is Doctor Who s05e02 an episode? %v\n", isEpisodeEpisode)

	// Output:
	// Is Nochnoy Dozor an episode? false
	// Is Doctor Who s05e02 an episode? true
}

func ExampleShow_SeasonEpisode() {
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()

	seasonNumber, episodeNumber, err := episode.SeasonEpisode()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("s%02de%02d\n", seasonNumber, episodeNumber)

	// Output:
	// s05e02
}

func ExampleShow_SeriesID() {
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()

	id, err := episode.SeriesID()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("id: %07d\n", id)

	// Output:
	// id: 0436992
}

func ExampleShow_SeriesTitle() {
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()

	title, err := episode.SeriesTitle()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("title: %s\n", title)

	// Output:
	// title: Doctor Who
}

func ExampleShow_SeriesYear() {
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()

	year, err := episode.SeriesYear()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("year: %s\n", year)

	// Output:
	// year: 2005
}
