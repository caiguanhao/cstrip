package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	DATA_FILE = "data/commitstrip.json"
)

type Document struct {
	URL     string    `json:"url"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	Image   string    `json:"image"`
	Content string    `json:"content"`
}

type CommitStrip struct {
	Content string `json:"content" binding:"required"`
}

func getDocuments() []Document {
	data, _ := ioutil.ReadFile(DATA_FILE)
	var documents []Document
	json.Unmarshal(data, &documents)
	return documents
}

func saveDocuments(documents []Document) {
	jsonFile, _ := json.MarshalIndent(documents, "", "  ")
	ioutil.WriteFile(DATA_FILE, append(jsonFile, '\n'), 0644)
}

func updateDocument(params martini.Params, cs CommitStrip, res http.ResponseWriter) {
	var index int = -1
	fmt.Sscanf(params["index"], "%d", &index)
	if index < 0 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	documents := getDocuments()
	if index > len(documents)-1 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	documents[index].Content = cs.Content
	saveDocuments(documents)
}

func main() {
	routes := martini.NewRouter()
	routes.Post("/update/:index", binding.Bind(CommitStrip{}), updateDocument)

	app := martini.New()
	staticOpts := martini.StaticOptions{SkipLogging: true}
	// app.Use(martini.Logger())
	app.Use(martini.Recovery())
	app.Use(martini.Static("public", staticOpts))
	app.Use(martini.Static("data", staticOpts))
	app.MapTo(routes, (*martini.Routes)(nil))
	app.Action(routes.Handle)
	app.Run()
}
