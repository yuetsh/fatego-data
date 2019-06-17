package main

import (
	"io"
	"net/http"
	"os"
)

// DownloadImage 下载图片
func DownloadImage(folder, filename, url string) {
	if url == "" {
		return
	}
	err := os.MkdirAll("images/"+folder, os.ModePerm)
	if err != nil {
		panic(err)
	}
	img, err := os.Create("images/" + folder + "/" + filename + ".png")
	if err != nil {
		panic(err)
	}
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(img, resp.Body)
	_ = img.Close()
}
