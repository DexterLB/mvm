package imdb

import (
	"fmt"
	"strings"

	_ "github.com/orchestrate-io/dvr"
)

func ExampleItem_ID() {
	movie := New(403358)
	defer movie.Free()

	fmt.Printf("id: %d\n", movie.ID())

	// Output:
	// id: 403358
}

func ExampleItem_Title() {
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

func ExampleItem_Title_episode() {
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

func ExampleItem_Title_series() {
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

func ExampleItem_Year() {
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

func ExampleItem_Year_episode() {
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

func ExampleItem_Year_series() {
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

func ExampleItem_OtherTitles() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	otherTitles, err := movie.OtherTitles()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	versions := []string{
		"UK",
		"Russia",
		"World-wide (English title)",
	}

	for _, version := range versions {
		fmt.Printf("%s: %s\n", version, otherTitles[version])
	}

	// Output:
	// UK: Night Watch
	// Russia: Ночной дозор
	// World-wide (English title): Night Watch
}

func ExampleItem_ReleaseDate() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	releaseDate, err := movie.ReleaseDate()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("release date: %s\n", releaseDate.Format("2006-01-02"))

	// Output:
	// release date: 2005-10-07
}

func ExampleItem_ReleaseDate_episode() {
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()

	releaseDate, err := episode.ReleaseDate()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("release date: %s\n", releaseDate.Format("2006-01-02"))

	// Output:
	// release date: 2010-04-10
}

func ExampleItem_Tagline() {
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

func ExampleItem_Languages_episode() {
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()

	languages, err := episode.Languages()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	languageNames := make([]string, len(languages))

	for i, language := range languages {
		languageNames[i] = language.String()
	}

	fmt.Printf("languages: %s\n", strings.Join(languageNames, ", "))

	// Output:
	// languages: en
}

func ExampleItem_Languages_movie() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	languages, err := movie.Languages()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	languageNames := make([]string, len(languages))

	for i, language := range languages {
		languageNames[i] = language.String()
	}

	fmt.Printf("languages: %s\n", strings.Join(languageNames, ", "))

	// Output:
	// languages: ru, de
}

func ExampleItem_Languages_series() {
	series := New(436992) // Doctor Who
	defer series.Free()

	languages, err := series.Languages()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	languageNames := make([]string, len(languages))

	for i, language := range languages {
		languageNames[i] = language.String()
	}

	fmt.Printf("languages: %s\n", strings.Join(languageNames, ", "))

	// Output:
	// languages: en
}

func ExampleItem_Duration() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	duration, err := movie.Duration()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("duration: %s\n", duration)

	// Output:
	// duration: 1h54m0s
}

func ExampleItem_Plot() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	shortPlot, err := movie.Plot()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("short plot: %s\n", shortPlot)

	// Output:
	// short plot: A fantasy-thriller set in present-day Moscow where the respective forces that control daytime and nighttime do battle.
}

func ExampleItem_PlotMedium() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	mediumPlot, err := movie.PlotMedium()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	sentences := strings.Split(mediumPlot, ". ")

	fmt.Printf("medium plot sentence 3:\n%s\n", sentences[2])

	// Output:
	// medium plot sentence 3:
	// Ever since, the forces of light govern the day while the night belongs to their dark opponents
}

func ExampleItem_PlotLong() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	longPlot, err := movie.PlotLong()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	sentences := strings.Split(longPlot, ". ")

	fmt.Printf("long plot sentence 3:\n%s\n", sentences[2])

	// Output:
	// long plot sentence 3:
	// They are known as the Others and have co-existed with humans for as long as humanity has existed
}

func ExampleItem_PosterURL() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	posterURL, err := movie.PosterURL()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("poster url: %s\n", posterURL)

	// Output:
	// poster url: http://ia.media-imdb.com/images/M/MV5BMjE0Nzk0NDkyOV5BMl5BanBnXkFtZTcwMjkzOTkyMQ@@.jpg
}

func ExampleItem_Rating() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	rating, err := movie.Rating()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("rating: %.2g\n", rating)

	// Output:
	// rating: 6.5
}

func ExampleItem_Votes() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	votes, err := movie.Votes()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("thousands of votes: %d\n", votes/1000)

	// Output:
	// thousands of votes: 46
}

func ExampleItem_Type() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()
	series := New(436992) // Doctor Who
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

func ExampleItem_SeasonEpisode() {
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

func ExampleItem_Series() {
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()

	series, err := episode.Series()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("id: %07d\n", series.ID())

	itemType, err := series.Type()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	fmt.Printf("type: %s\n", itemType)

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

func ExampleItem_Seasons_testnumbers() {
	series := New(1286039) // Stargate Universe
	defer series.Free()

	seasons, err := series.Seasons()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("seasons:\n")
	for i := range seasons {
		number, err := seasons[i].Number()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}

		fmt.Printf("%d\n", number)
	}

	// Output:
	// seasons:
	// 1
	// 2
}

func ExampleItem_Seasons_episodes() {
	series := New(1286039) // Stargate Universe
	defer series.Free()

	seasons, err := series.Seasons()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	if len(seasons) < 2 {
		fmt.Printf("season 2 is missing :(\n")
	}

	season2 := seasons[1]

	episodes, err := season2.Episodes()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	for i := range episodes {
		title, err := episodes[i].Title()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}

		// please note that i != episodeNumber, and possibly i+1 != episodeNumber
		seasonNumber, episodeNumber, err := episodes[i].SeasonEpisode()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return
		}

		fmt.Printf("s%02de%02d: %s\n", seasonNumber, episodeNumber, title)
	}

	// Output:
	// s02e01: Intervention
	// s02e02: Aftermath
	// s02e03: Awakening
	// s02e04: Pathogen
	// s02e05: Cloverdale
	// s02e06: Trial and Error
	// s02e07: The Greater Good
	// s02e08: Malice
	// s02e09: Visitation
	// s02e10: Resurgence
	// s02e11: Deliverance
	// s02e12: Twin Destinies
	// s02e13: Alliances
	// s02e14: Hope
	// s02e15: Seizure
	// s02e16: The Hunt
	// s02e17: Common Descent
	// s02e18: Epilogue
	// s02e19: Blockade
	// s02e20: Gauntlet
}

func ExampleItem_AllData_movie() {
	movie := New(403358) // Nochnoy Dozor (2004)
	defer movie.Free()

	data, err := movie.AllData()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("%s\n", data)

	// Output:
	// id: 403358
	// type: Movie
	// title: Nochnoy dozor
	// year: 2004
	// other titles:
	//  > Argentina -> Guardianes de la noche
	//  > Brazil -> Guardiões da Noite
	//  > Bulgaria (Bulgarian title) -> Нощна стража
	//  > Denmark -> Mørkets vogtere
	//  > Estonia -> Öine patrull
	//  > Finland -> Night Watch - yövahti
	//  > Finland (DVD title) -> Yövahti
	//  > Finland (Swedish title) -> Nattens väktare - nochnoi dozor
	//  > Finland (alternative title) -> Yövartija
	//  > France -> Night Watch
	//  > Georgia -> Gamis gushagi
	//  > Germany -> Wächter der Nacht - Nochnoi Dozor
	//  > Greece (transliterated ISO-LATIN-1 title) -> Oi fylakes tis nyhtas
	//  > Hungary -> Éjszakai őrség
	//  > Italy -> I guardiani della notte
	//  > Italy (DVD title) -> Night watch - I guardiani della notte
	//  > Latvia -> Nakts Sardze
	//  > Panama -> Guardianes de la noche
	//  > Peru -> Guardianes de la noche
	//  > Poland -> Straz nocna
	//  > Portugal -> Guardiões da Noite
	//  > Russia -> Ночной дозор
	//  > Serbia -> Noćna straža
	//  > Spain -> Guardianes de la noche
	//  > Sweden -> Nattens väktare
	//  > Turkey (Turkish title) -> Gece nöbeti
	//  > UK -> Night Watch
	//  > World-wide (English title) -> Night Watch
	// duration: 1h54m0s
	// short plot: A fantasy-thriller set in present-day Moscow where the respective forces that control daytime and nighttime do battle.
	// medium plot: Among normal humans live the "Others" possessing various supernatural powers...
	// long plot: THE SETTING: In the world that is modern Moscow, there exists a parallel realm known as the Gloom (kind of like the Astral Plane)...
	// poster url: http://ia.media-imdb.com/images/M/MV5BMjE0Nzk0NDkyOV5BMl5BanBnXkFtZTcwMjkzOTkyMQ@@.jpg
	// rating: 6.5
	// votes: 46k
	// languages: ru, de
	// release date: 2005-10-07
	// tagline: All That Stands Between Light And Darkness Is The Night Watch.
}

func ExampleItem_AllData_episode() {
	episode := New(1577257) // Doctor Who s05e02
	defer episode.Free()

	data, err := episode.AllData()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("%s\n", data)

	// Output:
	// id: 1577257
	// type: Episode
	// title: The Beast Below
	// year: 2010
	// other titles:
	//
	// duration: 42m0s
	// short plot: The Doctor takes Amy to the future inside Starship UK, which contains in addition to British explorers, an intimidating race known as the Smilers.
	// medium plot: The Doctor and Amy travel to a future time where all of the residents actually live in a orbiting spacecraft, Starship UK...
	// long plot: Starship UK is floating through space and we can see the words 'Yorkshire', 'Kent' and 'Surrey' visible on some of the buildings...
	// poster url: http://ia.media-imdb.com/images/M/MV5BNjY3MDI5OTE3N15BMl5BanBnXkFtZTcwMzA0MDU1NA@@.jpg
	// rating: 7.7
	// votes: 3k
	// languages: en
	// release date: 2010-04-10
	// season number: 5
	// episode number: 2
	// series id: 0436992
}

func ExampleItem_AllData_series() {
	series := New(1286039) // Stargate Universe
	defer series.Free()

	data, err := series.AllData()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf("%s\n", data)

	// Output:
	// id: 1286039
	// type: Series
	// title: SGU Stargate Universe
	// year: 2009
	// other titles:
	//  > France -> Stargate: Universe
	//  > Germany -> Stargate: Universe
	//  > Hungary -> Csillagkapu: Univerzum
	//  > Netherlands -> Stargate: Universe
	//  > Poland -> Gwiezdne Wrota: Wszechswiat
	//  > Russia -> Звёздные врата: Вселенная
	//  > Turkey (Turkish title) (new title) -> Yildiz Geçidi: Evren
	//  > USA (promotional title) -> SG.U Stargate Universe
	//  > World-wide (alternative title) (English title) -> Stargate: Universe
	// duration: 43m0s
	// short plot: Trapped on an Ancient spaceship billions of light years from home, a group of soldiers and civilians struggle to survive and find their way back to Earth.
	// medium plot: The Previously unknown purpose of the "Ninth Chevron" is revealed, and ends up taking a team to an Ancient ship "Destiny", a ship built millions of years ago by the Ancients, used to seed Distant galaxies with Stargates...
	// long plot: An attack on a secret off-world base by a rebel organisation has stranded the remaining survivors on an Ancient ship "Destiny", a large unmanned ship launched millions of years ago...
	// poster url: http://ia.media-imdb.com/images/M/MV5BOTEzNTY5NDY5M15BMl5BanBnXkFtZTcwMTY4MDQ3Mg@@.jpg
	// rating: 7.7
	// votes: 37k
	// languages: en
	// seasons: 1, 2
}
