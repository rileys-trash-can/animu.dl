package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"

	"crypto/md5"
	"encoding/hex"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

var c *colly.Collector

const URL = "https://animu.date"
const OutDir = "out/"

func main() {
	// ensure out dir
	err := os.Mkdir(OutDir, 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to ensure out dir: %s", err)
	}

	c = colly.NewCollector()

	c.OnHTML("img", DownloadImage)
	c.OnResponse(DownloadHandler)

	c.Visit(URL)

}

func DownloadHandler(res *colly.Response) {
	// non image
	if res.Headers.Get("Content-Type") == "text/html" {
		return
	}

	// build hash
	hashBytes := md5.Sum(res.Body)
	hashString := hex.EncodeToString(hashBytes[:4])

	log.Printf("Save %s (%d bytes): %s", hashString, len(res.Body), res.Request.URL.Path)

	dir, _ := filepath.Split(res.Request.URL.Path)
	err := EnsureDir(dir)
	if err != nil {
		log.Printf("Failed to create dir %s: %s", res.Request.URL.Path, err)
		return
	}

	err = res.Save(filepath.Join(OutDir,
		res.Request.URL.Path),
	)
	if err != nil {
		log.Printf("Failed to save: %s", err)
	}
}

func DownloadImage(h *colly.HTMLElement) {
	imgElem := h.DOM

	path := GetAttr(imgElem, "src")
	if path == "" {
		log.Printf("No src attribute!")
		return
	}

	// check if allready have
	_, err := os.Stat(filepath.Join(OutDir, path))
	if err == nil {
		log.Printf("Redundant download: %s", path)
		return
	}

	err = Download(path)
	if err != nil {
		log.Printf("Failed to download: %s", err)
		return
	}
}

func EnsureDir(path string) error {
	rebuildPath := OutDir

	paths := filepath.SplitList(path)
	for i := 0; i < len(paths); i++ {
		rebuildPath = filepath.Join(rebuildPath, paths[i])

		err := os.Mkdir(rebuildPath, 0750)

		if err != nil && !os.IsExist(err) {
			return err
		}
	}

	return nil
}

func Download(path string) error {
	url, err := url.Parse(URL)
	if err != nil {
		return err
	}

	url.Path = path

	return c.Visit(url.String())
}

func GetAttr(elem *goquery.Selection, attr string) string {
	if len(elem.Nodes) == 0 {
		return ""
	}

	for _, a := range elem.Nodes[0].Attr {
		if a.Key == attr {
			return a.Val
		}
	}

	return ""
}
