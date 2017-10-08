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

// GetRoutes returns the routes handled by this controller.
func (c *WebcamController) GetRoutes() []Route {
	return []Route{
		Route{"GET", "/", c.sendWebcamList},
		Route{"GET", "/:id", c.sendWebcam},
	}
}

func (c *WebcamController) sendWebcamList(w http.ResponseWriter, r *http.Request, p PathParams) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.Encode(c.webcams)
}

func (c *WebcamController) sendWebcam(w http.ResponseWriter, r *http.Request, p PathParams) {
	webcamID, err := strconv.Atoi(p["id"])
	if err != nil {
		notFound(w)
		return
	}

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
