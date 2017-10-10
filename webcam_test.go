package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var imageData = []byte{'T', 'e', 's', 't', '!'}

type testHandler struct{}

func (h testHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write(imageData)
}

func TestWebcamGetImage(t *testing.T) {
	server := httptest.NewServer(testHandler{})
	defer server.Close()

	addr := "http://" + server.Listener.Addr().String()

	webcam := &Webcam{
		1,
		"Les Paccots",
		addr,
		Coordinate{46.123, 6.66},
		"10",
		"3ms",
	}

	img, err := webcam.getImage(server.Client())
	if err != nil {
		t.Fail()
	}

	if !compare(img, imageData) {
		t.Fail()
	}
}

func compare(x, y []byte) bool {
	if len(x) != len(y) {
		return false
	}

	for i, v := range x {
		if v != y[i] {
			return false
		}
	}

	return true
}

func TestWebcamDurationParsing(t *testing.T) {
	webcam := &Webcam{
		1,
		"Les Paccots",
		"",
		Coordinate{46.123, 6.66},
		"10",
		"3ms",
	}

	if webcam.CrawlInterval() != 10*time.Second {
		t.Fail()
	}

	if webcam.MaxAge() != 3*time.Millisecond {
		t.Fail()
	}
}

func TestWebcamInvalidDurationParsing(t *testing.T) {
	webcam := &Webcam{
		1,
		"Les Paccots",
		"",
		Coordinate{46.123, 6.66},
		"hahaha",
		"3ms",
	}

	if webcam.CrawlInterval() != 0 {
		t.Fail()
	}
}
