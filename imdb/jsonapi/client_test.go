package jsonapi

import (
	"fmt"
	"net/http"

	"github.com/DexterLB/mvm/imdb"
)

func ExampleClient_Item() {
	client := &Client{
		HttpClient: http.DefaultClient,
		Address:    "http://localhost:8088",
	}

	data, err := client.Item(403358)
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
	// poster url: https://images-na.ssl-images-amazon.com/images/M/MV5BMjE0Nzk0NDkyOV5BMl5BanBnXkFtZTcwMjkzOTkyMQ@@.jpg
	// rating: 6.5
	// votes: 47k
	// languages: ru, de
	// release date: 2005-10-07
	// tagline: All That Stands Between Light And Darkness Is The Night Watch.
}

func ExampleClient_Search() {
	client := &Client{
		HttpClient: http.DefaultClient,
		Address:    "http://localhost:8088",
	}

	query := &imdb.SearchQuery{
		Query:    "Stalker",
		Year:     1979,
		Category: imdb.Movie,
		Exact:    true,
	}

	items, err := client.Search(query)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	for _, item := range items {
		fmt.Printf("[%s] %07d: %s (%d)\n", item.Type, item.ID, item.Title, item.Year)
	}

	// Output:
	// [Movie] 0079944: Stalker (1979)
}
