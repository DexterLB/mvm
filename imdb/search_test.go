package imdb

import "fmt"

func ExampleSearch() {
	query := &SearchQuery{
		Query: "Star Wars",
	}

	items, err := Search(query)

	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	for _, item := range items {
		title, err := item.Title()
		if err != nil {
			panic(err)
		}

		year, err := item.Year()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%07d: %s (%d)\n", item.ID(), title, year)
	}

	// Output:
	// 0076759: Star Wars (1977)
	// 0295630: Star Wars (1988)
	// 0251413: Star Wars (1983)
	// 2488496: Star Wars: Episode VII - The Force Awakens (2015)
	// 3748528: Rogue One: A Star Wars Story (2016)
	// 2527336: Star Wars: Episode VIII (2017)
	// 0458290: Star Wars: The Clone Wars (2008)
	// 2930604: Star Wars Rebels (2014)
	// 2527338: Star Wars: Episode IX (2019)
	// 1185834: Star Wars: The Clone Wars (2008)
}

func ExampleSearch_exact() {
	query := &SearchQuery{
		Query:    "Stalker",
		Year:     1979,
		Category: Movie,
		Exact:    true,
	}

	items, err := Search(query)

	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	for _, item := range items {
		title, err := item.Title()
		if err != nil {
			panic(err)
		}

		year, err := item.Year()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%07d: %s (%d)\n", item.ID(), title, year)
	}

	// Output:
	// 0079944: Stalker (1979)
}
