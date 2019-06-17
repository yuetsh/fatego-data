package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/exporter"
)

var File = "fatego_servants.json"

type Servant struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	Link  string            `json:"link"`
	Star  string            `json:"star"`
	Class string            `json:"class"`
	Image map[string]string `json:"image"`
}

func main() {
	file, err := os.Open(File)
	defer file.Close()
	if os.IsNotExist(err) {
		fetchServants()
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var servant Servant
		err = json.Unmarshal(scanner.Bytes(), &servant)
		for name, url := range servant.Image {
			DownloadImage(servant.ID+"_"+servant.Name, name, url)
			log.Println("Downloaded", name)
		}
	}
}

func fetchServants() {
	geziyor.NewGeziyor(geziyor.Options{
		StartURLs: []string{"http://wiki.joyme.com/fgo/%E8%8B%B1%E7%81%B5%E5%88%97%E8%A1%A8"},
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.80 Safari/537.36",
		Exporters: []geziyor.Exporter{exporter.JSONExporter{
			FileName: File,
		}},
		ParseFunc: func(r *geziyor.Response) {
			r.DocHTML.Find("#CardSelectTr > tbody").Children().Each(func(i int, s *goquery.Selection) {
				if i == 0 {
					return
				}
				// 基本信息
				id := s.Find("td:nth-child(1)")
				link := s.Find("td:nth-child(3) > a")
				star := s.Find("td:nth-child(4)")
				class := s.Find("td:nth-child(5)")
				url, _ := link.Attr("href")

				servant := Servant{
					ID:    strings.TrimSpace(id.Text()),
					Name:  link.Text(),
					Link:  "http://wiki.joyme.com" + url,
					Star:  strings.TrimSpace(star.Text()),
					Class: strings.TrimSpace(class.Text()),
					Image: make(map[string]string),
				}
				// 详细信息
				r.Geziyor.Get(servant.Link, func(r *geziyor.Response) {
					table := r.DocHTML.Find("#mw-content-text > div:nth-child(4) > div.col-md-4 > table")
					table.Find("img").Each(func(i int, s *goquery.Selection) {
						url, _ := s.Attr("src")
						name, _ := s.Attr("data-file-name")
						parsedName := strings.ReplaceAll(name, ".png", "")
						servant.Image[parsedName] = strings.ReplaceAll(url, "/dr/300__", "")
					})
				})
				r.Exports <- servant
			})
		},
	}).Start()
}
