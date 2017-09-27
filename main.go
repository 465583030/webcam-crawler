package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func loadWebcams() []Webcam {
	file, err := os.Open("webcams.json")
	if err != nil {
		panic("Could not read configuration file")
	}

	defer file.Close()

	webcams := []Webcam{}
	decoder := json.NewDecoder(file)
	decoder.Decode(&webcams)

	return webcams
}

func main() {
	controller := WebcamController{}
	controller.SetWebcams(loadWebcams())

	http.Handle("/", controller)
	http.ListenAndServe(":8080", nil)
}
