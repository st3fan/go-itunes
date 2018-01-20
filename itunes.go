package itunes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/mitchellh/mapstructure"
)

type Content interface {
	ContentType() string
	ContentKind() string
}

type Software struct {
	WrapperType                        string  `json:"wrapperType"`
	Kind                               string  `json:"kind"`
	ArtistName                         string  `json:"artistName"`
	ArtistID                           int     `json:"artistId"`
	TrackName                          string  `json:"trackName"`
	TrackCensoredName                  string  `json:"trackCensoredName"`
	TrackID                            int     `json:"trackId"`
	AverageUserRating                  float64 `json:"averageUserRating"`
	UserRatingCount                    int     `json:"userRatingCount"`
	AverageUserRatingForCurrentVersion float64 `json:"averageUserRatingForCurrentVersion"`
	UserRatingCountForCurrentVersion   int     `json:"userRatingCountForCurrentVersion"`
	Version                            string  `json:"string"`
}

func (o *Software) ContentType() string {
	return o.WrapperType
}

func (o *Software) ContentKind() string {
	return o.ContentKind()
}

type Podcast struct {
	WrapperType string `json:"wrapperType"`
	Kind        string `json:"kind"`
	FeedURL     string `json:"feedUrl"`
}

func (o *Podcast) ContentType() string {
	return o.WrapperType
}

func (o *Podcast) ContentKind() string {
	return o.ContentKind()
}

type Song struct {
	WrapperType string `json:"wrapperType"`
	Kind        string `json:"kind"`
}

func (o *Song) ContentType() string {
	return o.WrapperType
}

func (o *Song) ContentKind() string {
	return o.ContentKind()
}

type Artist struct {
	WrapperType string `json:"wrapperType"`
	ArtistName  string `json:"artistName"`
}

func (o *Artist) ContentType() string {
	return o.WrapperType
}

func (o *Artist) ContentKind() string {
	return ""
}

type SearchOptions struct {
	Term      string `url:"term"`
	Country   string `url:"country,omitempty"`
	Media     string `url:"media,omitempty"`
	Entity    string `url:"entity,omitempty"`
	Attribute string `url:"attribute,omitempty"`
	Limit     int    `url:"limit,omitempty"`
	Language  string `url:"language,omitempty"`
	Explicit  bool   `url:"explicit,omitempty"`
}

func Country(country string) func(*SearchOptions) error {
	return func(so *SearchOptions) error {
		so.Country = country
		return nil
	}
}

func Media(media string) func(*SearchOptions) error {
	return func(so *SearchOptions) error {
		so.Media = media
		return nil
	}
}

func Entity(entity string) func(*SearchOptions) error {
	return func(so *SearchOptions) error {
		so.Entity = entity
		return nil
	}
}

func Attribute(attribute string) func(*SearchOptions) error {
	return func(so *SearchOptions) error {
		so.Attribute = attribute
		return nil
	}
}

func Explicit(explicit bool) func(*SearchOptions) error {
	return func(so *SearchOptions) error {
		so.Explicit = explicit
		return nil
	}
}

func Limit(limit int) func(*SearchOptions) error {
	return func(so *SearchOptions) error {
		so.Limit = limit
		return nil
	}
}

func Language(language string) func(*SearchOptions) error {
	return func(so *SearchOptions) error {
		so.Language = language
		return nil
	}
}

type searchResults struct {
	ResultCount int                      `json:"resultCount"`
	Results     []map[string]interface{} `json:"results"`
}

func parseSearchResults(results searchResults) ([]Content, error) {
	var content []Content
	for _, object := range results.Results {
		if wrapperType, ok := object["wrapperType"]; ok {
			switch wrapperType {
			case "track":
				if kind, ok := object["kind"]; ok {
					switch kind {
					case "song":
						var song Song
						if err := mapstructure.Decode(object, &song); err != nil {
							return nil, err
						}
						content = append(content, &song)
					case "podcast":
						var podcast Podcast
						if err := mapstructure.Decode(object, &podcast); err != nil {
							return nil, err
						}
						content = append(content, &podcast)
					}
				}
			case "collection":
				break
			case "artist":
				var artist Artist
				if err := mapstructure.Decode(object, &artist); err != nil {
					return nil, err
				}
				content = append(content, &artist)
			}
		}
	}

	return content, nil
}

func queryServer(url string) ([]byte, error) {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid status <%s>", res.Status)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, err
	}

	return body, nil
}

func Search(term string, options ...func(*SearchOptions) error) ([]Content, error) {
	searchOptions := SearchOptions{Term: term}
	for _, option := range options {
		if err := option(&searchOptions); err != nil {
			return nil, err
		}
	}

	values, err := query.Values(searchOptions)
	if err != nil {
		return nil, err
	}

	body, err := queryServer(fmt.Sprintf("https://itunes.apple.com/search?" + values.Encode()))
	if err != nil {
		return nil, err
	}

	results := searchResults{}
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	return parseSearchResults(results)
}

type LookupOptions struct {
	ID int `url:"id"`
}

func Lookup(id int, options ...func(*LookupOptions) error) ([]Content, error) {
	lookupOptions := LookupOptions{ID: id}
	for _, option := range options {
		if err := option(&lookupOptions); err != nil {
			return nil, err
		}
	}

	values, err := query.Values(lookupOptions)
	if err != nil {
		return nil, err
	}

	body, err := queryServer(fmt.Sprintf("https://itunes.apple.com/lookup?" + values.Encode()))
	if err != nil {
		return nil, err
	}

	results := searchResults{}
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	return parseSearchResults(results)
}
