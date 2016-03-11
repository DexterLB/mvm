package imdb

import (
	"fmt"
	"strings"
)

func ExampleShow_Title() {
	movie := New("0403358") // Nochnoy Dozor (2004)

	fmt.Printf("title: %s\n", movie.Title())

	// Output:
	// title: Nochnoy dozor
}

func ExampleShow_OtherTitles() {
	movie := New("0403358") // Nochnoy Dozor (2004)

	fmt.Printf("other titles: %s\n", strings.Join(movie.OtherTitles(), ", "))

	// Output:
	// other titles: Night Watch
}

func ExampleShow_ReleaseDate() {
	movie := New("0403358") // Nochnoy Dozor (2004)

	fmt.Printf("release date: %s\n", movie.ReleaseDate())

	// Output:
	// release date: 2006-03-03
}

func ExampleShow_Tagline() {
	movie := New("0403358") // Nochnoy Dozor (2004)

	fmt.Printf("tagline: %s\n", movie.Tagline())

	// Output:
	// tagline: All That Stands Between Light And Darkness Is The Night Watch.
}

func ExampleShow_Duration() {
	movie := New("0403358") // Nochnoy Dozor (2004)

	fmt.Printf("duration: %d\n", movie.Duration())

	// Output:
	// duration: 1h 54m
}

func ExampleShow_Plot() {
	movie := New("0403358") // Nochnoy Dozor (2004)

	fmt.Printf("short plot: %s\n", movie.Plot(0))

	// Output:
	// short plot: A fantasy-thriller set in present-day Moscow where the respective forces that control daytime and nighttime do battle.
}

func ExampleShow_PosterURL() {
	movie := New("0403358") // Nochnoy Dozor (2004)

	fmt.Printf("poster url: %s\n", movie.PosterURL())

	// Output:
	// poster url: http://ia.media-imdb.com/images/M/MV5BMjE0Nzk0NDkyOV5BMl5BanBnXkFtZTcwMjkzOTkyMQ@@._V1_SX640_SY720_.jpg
}

func ExampleShow_Rating() {
	movie := New("0403358") // Nochnoy Dozor (2004)

	fmt.Printf("rating: %.1g\n", movie.Rating())

	// Output:
	// rating: 6.5
}

func ExampleShow_Votes() {
	movie := New("0403358") // Nochnoy Dozor (2004)

	fmt.Printf("votes: %d\n", movie.Votes())

	// Output:
	// votes: 46450
}
