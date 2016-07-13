package imdb

import (
	"testing"

	_ "github.com/orchestrate-io/dvr"
	"github.com/stretchr/testify/assert"
)

func TestRatingEdgeCases(t *testing.T) {
	episode := New(5128652) // Magicians s01e06
	defer episode.Free()

	rating, err := episode.Rating()
	if err != nil {
		t.Fatal(err)
	}

	votes, err := episode.Votes()
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.InDelta(8.4, rating, 0.1)
	assert.Equal(3, votes/100)
}

func TestRatingEdgeCases_AllData(t *testing.T) {
	episode := New(5128652) // Magicians s01e06
	defer episode.Free()

	data, err := episode.AllData()
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.InDelta(8.4, data.Rating, 0.1)
	assert.Equal(3, data.Votes/100)
}
