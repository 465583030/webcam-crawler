package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
)

// WebcamController struct contains the webcam data and provides methods to handle HTTP requests.
type WebcamController struct {
	client      http.Client
	webcams     []Webcam
	storagePath string
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
		Route{"GET", "/:id/hist", c.sendHist},
		Route{"GET", "/:id/hist/:name", c.sendHistWebcam},
	}
}

func (c *WebcamController) sendWebcamList(w http.ResponseWriter, r *http.Request, p PathParams) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.Encode(c.webcams)
}

func (c *WebcamController) sendWebcam(w http.ResponseWriter, r *http.Request, p PathParams) {
	webcam := c.getWebcamOrNotFound(p["id"], w)
	if webcam == nil {
		return
	}

	imageBytes, err := webcam.getImage(&c.client)
	if err != nil {
		proxyFailure(w)
		return
	}

	w.Header().Add("Content-Type", "image/jpeg")
	w.Write(imageBytes)
}

func (c *WebcamController) sendHist(w http.ResponseWriter, r *http.Request, p PathParams) {
	w.Header().Set("Content-Type", "application/json")

	hist := []string{}
	encoder := json.NewEncoder(w)

	if c.getWebcamOrNotFound(p["id"], w) == nil {
		return
	}

	path := filepath.Join(c.storagePath, p["id"])
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("Unable to read directory %s: %s\n", path, err)
		encoder.Encode(hist)
		return
	}

	for _, f := range files {
		hist = append(hist, f.Name())
	}

	encoder.Encode(hist)
}

func (c *WebcamController) sendHistWebcam(w http.ResponseWriter, r *http.Request, p PathParams) {
	if c.getWebcamOrNotFound(p["id"], w) == nil {
		return
	}

	path := filepath.Join(c.storagePath, p["id"], p["name"])
	imageBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Unable to read file: %s: %s", path, err)
		notFound(w)
		return
	}

	w.Write(imageBytes)
}

func (c *WebcamController) getWebcamOrNotFound(id string, w http.ResponseWriter) *Webcam {
	webcamID, err := strconv.Atoi(id)
	if err != nil {
		notFound(w)
		return nil
	}

	for _, w := range c.webcams {
		if w.ID == webcamID {
			return &w
		}
	}

	notFound(w)
	return nil
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func proxyFailure(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadGateway)
}
