package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// A Crawler crawls webcam images at the interval specified
// in Webcam structs and saves them in a local folder.
type Crawler struct {
	webcams     []Webcam
	client      *http.Client
	storagePath string
	stopChans   []chan struct{}
	format      string
}

// NewCralwer creates a new crawler given a list of Webcams and a path to a folder to save images.
func NewCralwer(webcams []Webcam, storagePath string) *Crawler {
	return &Crawler{
		webcams,
		&http.Client{},
		storagePath,
		make([]chan struct{}, 0),
		time.RFC3339,
	}
}

// Start the crawler.
func (c *Crawler) Start() {
	for _, w := range c.webcams {
		c.scheduleCrawl(w)
	}
}

// Stop the crawler.
func (c *Crawler) Stop() {
	for _, ch := range c.stopChans {
		close(ch)
	}
}

func (c *Crawler) scheduleCrawl(w Webcam) {
	if w.CrawlInterval() <= 0 {
		// Crawling is disabled for this webcam
		return
	}

	stopChan := make(chan struct{})
	c.stopChans = append(c.stopChans, stopChan)

	ticker := time.NewTicker(w.CrawlInterval())

	go func() {
		for {
			select {
			case <-ticker.C:
				c.crawl(w)
			case <-stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (c *Crawler) crawl(w Webcam) {
	now := time.Now()

	image, err := w.getImage(c.client)
	if err != nil {
		fmt.Printf("Could not get image for webcam %s: %s\n", w.Name, err) // TODO: replace with proper logging
		return
	}

	dirname := filepath.Join(c.storagePath, strconv.Itoa(w.ID))
	os.MkdirAll(dirname, os.ModePerm)
	filename := filepath.Join(dirname, now.Format(c.format)+".jpg")

	err = ioutil.WriteFile(filename, image, 0644)
	if err != nil {
		fmt.Printf("Could not write file %s: %s\n", filename, err)
		return
	}

	c.cleanupDir(dirname, w.MaxAge())
}

func (c *Crawler) cleanupDir(dirname string, maxAge time.Duration) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Printf("Could not read directory %s: %s\n", dirname, err)
	}

	for _, f := range files {
		creationTime, err := c.timeFromName(f.Name())
		if err != nil {
			continue
		}

		now := time.Now()

		if now.Sub(creationTime) > maxAge {
			filename := filepath.Join(dirname, f.Name())
			err := os.Remove(filename)
			if err != nil {
				fmt.Printf("Could not remove file %s: %s", filename, err)
			}
		}
	}
}

func (c *Crawler) timeFromName(name string) (time.Time, error) {
	var extension = filepath.Ext(name)
	var dateString = name[0 : len(name)-len(extension)]

	creation, err := time.Parse(c.format, dateString)
	if err != nil {
		return time.Time{}, err
	}

	return creation, nil
}
