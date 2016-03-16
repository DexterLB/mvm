package imdb

import (
	"fmt"
	"strings"
)

func ExampleShow_ID() {
	movie := New(403358)
	fmt.Printf("id: %d\n", movie.ID())

	// Output:
	// id: 403358
}

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

func ExampleShow_Title_series() {
	series := New(436992) // Doctor Who
	defer series.Free()

	title, err := series.Title()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("title: %s\n", title)

	// Output:
	// title: Doctor Who
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

func ExampleShow_Year_episode() {
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()

	year, err := episode.Year()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("year: %d\n", year)

	// Output:
	// year: 2010
}

func ExampleShow_Year_series() {
	series := New(436992) // Doctor Who
	defer series.Free()

	year, err := series.Year()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("year: %d\n", year)

	// Output:
	// year: 2005
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

func ExampleShow_ReleaseDate_episode() {
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()

	releaseDate, err := episode.ReleaseDate()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("release date: %s\n", releaseDate)

	// Output:
	// release date: 2010-04-10
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

func ExampleShow_Type() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()
	series := New(436992)
	defer series.Free()

	movieType, err := movie.Type()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	episodeType, err := episode.Type()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	seriesType, err := series.Type()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("Nochnoy Dozor's type is %s\n", movieType)
	fmt.Printf("Doctor Who s05e02's type is %s\n", episodeType)
	fmt.Printf("Doctor Who's type is %s\n", seriesType)

	// Output:
	// Nochnoy Dozor's type is Movie
	// Doctor Who s05e02's type is Episode
	// Doctor Who's type is Series
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

func ExampleShow_Series() {
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()

	series, err := episode.Series()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("id: %07d\n", series.ID())

	showType, err := series.Type()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	fmt.Printf("type: %s\n", showType)

	title, err := series.Title()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	fmt.Printf("title: %s\n", title)

	// Output:
	// id: 0436992
	// type: Series
	// title: Doctor Who
}
