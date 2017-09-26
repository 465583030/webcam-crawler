package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var allWebcams = []Webcam{}

func loadWebcams() {
	file, err := os.Open("webcams.json")
	if err != nil {
		panic("Could not read configuration file")
	}

	decoder := json.NewDecoder(file)
	decoder.Decode(&allWebcams)
}

func getWebcam(name string) *Webcam {
	for _, w := range allWebcams {
		if w.Name == name {
			return &w
		}
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	webcamName := r.URL.RequestURI()[1:]
	webcam := getWebcam(webcamName)

	if webcam == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	imageBytes, err := getWebcam(webcamName).getImage()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not get image")
		return
	}

	w.Header().Add("Content-Type", "img/jpeg")
	w.Write(imageBytes)
}

func main() {
	loadWebcams()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
