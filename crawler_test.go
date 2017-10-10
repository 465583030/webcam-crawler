package main

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCrawler(t *testing.T) {
	server := httptest.NewServer(testHandler{})
	defer server.Close()

	addr := "http://" + server.Listener.Addr().String()

	webcams := []Webcam{
		Webcam{
			1,
			"Les Paccots",
			addr,
			Coordinate{46.123, 6.66},
			"3ms",
			"5ms",
		},
		Webcam{
			2,
			"La Fouly",
			addr,
			Coordinate{46.123, 6.66},
			"0",
			"5ms",
		}}

	storagePath := "test-storage"
	os.Mkdir(storagePath, os.ModePerm)

	c := NewCralwer(webcams, storagePath)
	c.client = server.Client()
	c.format = time.RFC3339Nano

	c.Start()
	time.Sleep(10 * time.Millisecond)
	c.Stop()
	time.Sleep(1 * time.Millisecond)

	files, err := ioutil.ReadDir(filepath.Join(storagePath, "1"))
	if err != nil {
		t.Fail()
	}

	if len(files) != 2 {
		t.Fail()
	}

	if _, err := os.Stat(filepath.Join(storagePath, "2")); !os.IsNotExist(err) {
		t.Fail()
	}

	os.RemoveAll(storagePath)
}
