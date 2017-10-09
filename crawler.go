package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// A Crawler crawls webcam images at the interval specified
// in Webcam structs and saves them in a local folder.
type Crawler struct {
	webcams     []Webcam
	storagePath string
	stopChans   []chan struct{}
}

// NewCralwer creates a new crawler given a list of Webcams and a path to a folder to save images.
func NewCralwer(webcams []Webcam, storagePath string) *Crawler {
	return &Crawler{
		webcams,
		storagePath,
		make([]chan struct{}, 0),
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
	image, err := w.getImage()
	if err != nil {
		fmt.Printf("Could not get image for webcam %s: %s\n", w.Name, err) // TODO: replace with proper logging
		return
	}

	dirname := filepath.Join(c.storagePath, strconv.Itoa(w.ID))
	os.MkdirAll(dirname, os.ModeDir|os.ModePerm)

	now := time.Now()
	filename := filepath.Join(dirname, now.Format(time.RFC3339)+".jpg")

	err = ioutil.WriteFile(filename, image, 0644)
	if err != nil {
		fmt.Printf("Could not write file %s: %s\n", filename, err)
		return
	}

	cleanupDir(dirname, w.MaxAge())
}

func cleanupDir(dirname string, maxAge time.Duration) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Printf("Could not read directory %s: %s\n", dirname, err)
	}

	for _, f := range files {
		creationTime, err := timeFromName(f.Name())
		if err != nil {
			continue
		}

		if time.Now().Sub(creationTime) > maxAge {
			filename := filepath.Join(dirname, f.Name())
			err := os.Remove(filename)
			if err != nil {
				fmt.Printf("Could not remove file %s: %s", filename, err)
			}
		}
	}
}

func timeFromName(name string) (time.Time, error) {
	nameParts := strings.Split(name, ".")
	if len(nameParts) != 2 {
		return time.Time{}, errors.New("Cannot extract time from name: " + name)
	}

	creation, err := time.Parse(time.RFC3339, nameParts[0])
	if err != nil {
		return time.Time{}, err
	}

	return creation, nil
}
