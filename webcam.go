package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Coordinate represents a geographic coordinate with latitude and longitude.
type Coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Webcam struct contains information concerning a webcam such
// as its name and the URL at which the webcam image can be retrieved.
type Webcam struct {
	ID                  int        `json:"id"`
	Name                string     `json:"name"`
	URL                 string     `json:"URL"`
	Position            Coordinate `json:"position"`
	CrawlIntervalString string     `json:"crawlInterval"`
	MaxAgeString        string     `json:"maxAge"`
}

// CrawlInterval returns the Duration between two image fetches.
func (w *Webcam) CrawlInterval() time.Duration {
	return myParseDuration(w.CrawlIntervalString)
}

// MaxAge returns the maximum Duration an image should be stored.
func (w *Webcam) MaxAge() time.Duration {
	return myParseDuration(w.MaxAgeString)
}

func (w *Webcam) getImage(client *http.Client) ([]byte, error) {
	r, err := client.Get(w.URL)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func myParseDuration(duration string) time.Duration {
	v, err := strconv.Atoi(duration)
	if err == nil {
		return time.Duration(v) * time.Second
	}

	d, err := time.ParseDuration(duration)
	if err != nil {
		return 0
	}

	return d
}
