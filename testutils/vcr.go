package testutils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

// RecordHTTP wraps tests and records all http requests made with the default
// http transport in a file with the given name. If the file exists,
// requests are replayed.
func RecordHTTP(m *testing.M, filename string) int {
	t := http.DefaultTransport

	rec, err := recorder.New(filename)
	if err != nil {
		panic(err)
	}

	rec.SetMatcher(func(rec *http.Request, i cassette.Request) bool {
		if !cassette.DefaultMatcher(rec, i) {
			return false
		}

		if rec.Body == nil {
			return true
		}

		var b bytes.Buffer
		if _, err := b.ReadFrom(rec.Body); err != nil {
			return false
		}
		rec.Body = ioutil.NopCloser(&b)

		return b.String() == "" || b.String() == i.Body
	})

	defer func() {
		http.DefaultTransport = t
		rec.Stop()
	}()

	http.DefaultTransport = rec

	return m.Run()
}
