package main

import (
	"encoding/json"
	"errors"
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

func (c *WebcamController) sendWebcamList(w http.ResponseWriter, r *http.Request, p PathParams) error {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.Encode(c.webcams)

	return nil
}

func (c *WebcamController) sendWebcam(w http.ResponseWriter, r *http.Request, p PathParams) error {
	webcam, err := c.getWebcam(p["id"], w)
	if err != nil {
		return err
	}

	imageBytes, err := webcam.getImage(&c.client)
	if err != nil {
		return StatusError{http.StatusBadGateway, err}
	}

	w.Header().Add("Content-Type", "image/jpeg")
	w.Write(imageBytes)

	return nil
}

func (c *WebcamController) sendHist(w http.ResponseWriter, r *http.Request, p PathParams) error {
	w.Header().Set("Content-Type", "application/json")

	hist := []string{}
	encoder := json.NewEncoder(w)

	// Check that webcam exists
	if _, err := c.getWebcam(p["id"], w); err != nil {
		return err
	}

	path := filepath.Join(c.storagePath, p["id"])
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("Unable to read directory %s: %s\n", path, err)
		encoder.Encode(hist)
		return nil
	}

	for _, f := range files {
		hist = append(hist, f.Name())
	}

	encoder.Encode(hist)
	return nil
}

func (c *WebcamController) sendHistWebcam(w http.ResponseWriter, r *http.Request, p PathParams) error {
	// Check that webcam exists
	if _, err := c.getWebcam(p["id"], w); err != nil {
		return err
	}

	path := filepath.Join(c.storagePath, p["id"], p["name"])
	imageBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return StatusError{http.StatusNotFound, errors.New("Unable to read file " + path + " : " + err.Error())}
	}

	w.Write(imageBytes)
	return nil
}

func (c *WebcamController) getWebcam(id string, w http.ResponseWriter) (*Webcam, error) {
	webcamID, err := strconv.Atoi(id)
	if err != nil {
		return nil, StatusError{http.StatusNotFound, errors.New("Cuuld not convert " + id + " to webcam id")}
	}

	for _, w := range c.webcams {
		if w.ID == webcamID {
			return &w, nil
		}
	}

	return nil, StatusError{http.StatusNotFound, errors.New("Could not find webcam with id " + id)}
}
