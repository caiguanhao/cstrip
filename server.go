package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"io/ioutil"
	"net/http"
	"os"
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

func updateDocuments(
	params martini.Params,
	cs CommitStrip,
	res http.ResponseWriter,
	req *http.Request,
) {
	reqUser := req.Header.Get("USERNAME")
	reqPass := req.Header.Get("PASSWORD")
	if reqUser != username || reqPass != password {
		res.WriteHeader(http.StatusUnauthorized)
	}

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

var username, password string

func main() {
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")

	if martini.Env == martini.Prod && (username == "" || password == "") {
		fmt.Fprintln(os.Stderr, "Please provide username and password!")
		os.Exit(1)
	}

	routes := martini.NewRouter()
	routes.Post("/update/:index", binding.Bind(CommitStrip{}), updateDocuments)

	app := martini.New()
	staticOpts := martini.StaticOptions{SkipLogging: true}
	// app.Use(martini.Logger())
	app.Use(martini.Recovery())
	if martini.Env == martini.Prod {
		app.Use(martini.Static("dist", staticOpts))
	} else {
		app.Use(martini.Static("public", staticOpts))
		app.Use(martini.Static("data", staticOpts))
	}
	app.MapTo(routes, (*martini.Routes)(nil))
	app.Action(routes.Handle)
	app.Run()
}
