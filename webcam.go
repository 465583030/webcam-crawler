package main

import (
	"io/ioutil"
	"net/http"
)

// Coordinate represents a geographic coordinate with latitude and longitude.
type Coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Webcam struct contains information concerning a webcam such
// as its name and the URL at which the webcam image can be retrieved.
type Webcam struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	URL      string     `json:"URL"`
	Position Coordinate `json:"position"`
}

func (w *Webcam) getImage() ([]byte, error) {
	r, err := http.Get(w.URL)

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
