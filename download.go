package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	DATA_FILE = "data/commitstrip.json"
)

type Document struct {
	Image string `json:"image"`
}

func download(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	contentLength := resp.Header.Get("Content-Length")
	var size int64
	fmt.Sscanf(contentLength, "%d", &size)
	filename := strings.Replace("/"+url, "/http://www.commitstrip.com/wp-content/uploads/", "", -1)
	if strings.HasPrefix(filename, "/") {
		panic(filename)
	}
	dirname := filepath.Join("dist", "strips", filepath.Dir(filename))
	os.MkdirAll(dirname, 0755)
	path := filepath.Join(dirname, filepath.Base(filename))
	info, err := os.Stat(path)
	if err == nil {
		if info.Size() == size {
			fmt.Printf("Already downloaded %s.\n", filename)
			return
		}
	}
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	written, err := io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received %d bytes of %s.\n", written, filename)
}

func getImages() []string {
	data, _ := ioutil.ReadFile(DATA_FILE)
	var documents []Document
	json.Unmarshal(data, &documents)

	var images []string
	for _, doc := range documents {
		imgs := strings.Split(doc.Image, "\n")
		for _, img := range imgs {
			images = append(images, img)
		}
	}
	return images
}

func batch(images *[]string, start, length int) {
	promise := make(chan int, length)
	for i := 0; i < length; i++ {
		go func(index int) {
			index += start
			if index < len(*images) {
				download((*images)[index])
			}
			promise <- index
		}(i)
	}
	for i := 0; i < length; i++ {
		<-promise
	}
}

func main() {
	os.MkdirAll(filepath.Join("dist", "strips"), 0755)
	images := getImages()
	imagesLen := len(images)
	imagesPerBatch := 10
	batches := int(math.Ceil(float64(imagesLen) / float64(imagesPerBatch)))
	for i := 0; i < batches; i++ {
		batch(&images, i*imagesPerBatch, imagesPerBatch)
	}
}
