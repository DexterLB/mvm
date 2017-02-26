package testutils

import (
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
)

func RecordHTTP(m *testing.M, cassette string) int {
	t := http.DefaultTransport

	rec, err := recorder.New(cassette)
	if err != nil {
		panic(err)
	}

	defer func() {
		http.DefaultTransport = t
		rec.Stop()
	}()

	http.DefaultTransport = rec

	return m.Run()
}
