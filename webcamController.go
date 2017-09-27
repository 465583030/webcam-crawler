package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// WebcamController struct contains the webcam data and provides methods to handle HTTP requests.
type WebcamController struct {
	webcams []Webcam
}

// SetWebcams sets the list of webcams that the controller can display.
func (c *WebcamController) SetWebcams(webcams []Webcam) {
	c.webcams = webcams
}

func (c WebcamController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.RequestURI() == "/" {
		c.sendWebcamList(w, r)
		return
	}

	webcamID, err := strconv.Atoi(r.URL.RequestURI()[1:])
	if err != nil {
		notFound(w)
		return
	}

	c.sendWebcam(webcamID, w, r)
}

func (c *WebcamController) sendWebcamList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.Encode(c.webcams)
}

func (c *WebcamController) sendWebcam(webcamID int, w http.ResponseWriter, r *http.Request) {
	webcam := c.getWebcam(webcamID)
	if webcam == nil {
		notFound(w)
		return
	}

	imageBytes, err := webcam.getImage()
	if err != nil {
		proxyFailure(w)
		return
	}

	w.Header().Add("Content-Type", "image/jpeg")
	w.Write(imageBytes)
}

func (c *WebcamController) getWebcam(id int) *Webcam {
	for _, w := range c.webcams {
		if w.ID == id {
			return &w
		}
	}
	return nil
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func proxyFailure(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadGateway)
}
