package itunes_test

import (
	"testing"

	"github.com/st3fan/itunes"
)

func Test_Lookup(t *testing.T) {
	results, err := itunes.Lookup(528458508)
	if err != nil {
		t.Fatal("Lookup failed: ", err)
	}

	if len(results) != 1 {
		t.Fatal("Expected 1 result; got <%d>", len(results))
	}

	podcast, ok := results[0].(*itunes.Podcast)
	if !ok {
		t.Fatal("Expected a podcast result; got <%v>", results[0])
	}

	if podcast.FeedURL != "https://daringfireball.net/thetalkshow/rss" {
		t.Fatal("Expected feed URL")
	}
}

func Test_Search(t *testing.T) {
	_, err := itunes.Search("jack johnson")
	if err != nil {
		t.Fatal("Search failed: ", err)
	}
}

func Test_SearchWithOptions(t *testing.T) {
	_, err := itunes.Search("Love", itunes.Media("ebook"))
	if err != nil {
		t.Fatal("Search failed: ", err)
	}
}
