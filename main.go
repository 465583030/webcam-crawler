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

func defaultHandler(w http.ResponseWriter, r *http.Request, p PathParams) {
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	controller := &WebcamController{}
	controller.SetWebcams(loadWebcams())

	router := NewRouter(defaultHandler)
	router.Mount("/", controller)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
