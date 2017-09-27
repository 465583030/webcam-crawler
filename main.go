package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var allWebcams = []Webcam{}

func loadWebcams() {
	file, err := os.Open("webcams.json")
	if err != nil {
		panic("Could not read configuration file")
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.Decode(&allWebcams)
}

func getWebcam(id int) *Webcam {
	for _, w := range allWebcams {
		if w.ID == id {
			return &w
		}
	}
	return nil
}

func sendWebcamList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.Encode(allWebcams)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.RequestURI() == "/" {
		sendWebcamList(w, r)
		return
	}

	webcamID, err := strconv.Atoi(r.URL.RequestURI()[1:])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	webcam := getWebcam(webcamID)
	if webcam == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	imageBytes, err := getWebcam(webcamID).getImage()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not get image")
		return
	}

	w.Header().Add("Content-Type", "image/jpeg")
	w.Write(imageBytes)
}

func main() {
	loadWebcams()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
