package main

import (
	"encoding/json"
	"errors"
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

func startCrawler(webcams []Webcam) {
	crawler := NewCralwer(webcams, "hist")
	crawler.Start()
}

func startWebServer(webcams []Webcam) {
	controller := &WebcamController{
		storagePath: "hist",
	}
	controller.SetWebcams(webcams)

	router := NewRouter(defaultHandler)
	router.Mount("/webcam", controller)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func defaultHandler(w http.ResponseWriter, r *http.Request, p PathParams) error {
	return StatusError{404, errors.New("Page not found at " + r.URL.Path)}
}

func main() {
	webcams := loadWebcams()

	startCrawler(webcams)
	startWebServer(webcams)
}
