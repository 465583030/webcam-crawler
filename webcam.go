package main

import (
	"io/ioutil"
	"net/http"
)

// Webcam struct contains information concerning a webcam such
// as its name and the URL at which the webcam image can be retrieved.
type Webcam struct {
	Name string
	URL  string
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

// AllWebcams contains all the webcams
var AllWebcams []Webcam = []Webcam{
	Webcam{"paccots", "http://manage.mycity.travel/proxy/webcam/51.jpg"},
	Webcam{"arpette", "http://panorama.simwatch.ch/image.php?cname=Arpette"}}
