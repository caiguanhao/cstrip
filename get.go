package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

type ByDate []Document

func (a ByDate) Len() int {
	return len(a)
}

func (a ByDate) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByDate) Less(i, j int) bool {
	return a[j].Date.Before(a[i].Date)
}

var (
	COMMITSTRIP_HOME_URL = "http://www.commitstrip.com/en/"
	COMMITSTRIP_PAGE_URL = "http://www.commitstrip.com/en/page/%d/"
)

func debug(a ...interface{}) {
	fmt.Fprint(os.Stderr, "DEBUG ")
	fmt.Fprintln(os.Stderr, a...)
}

func parseMonth(month string) time.Month {
	switch month {
	default:
		return time.January
	case "February":
		return time.February
	case "March":
		return time.March
	case "April":
		return time.April
	case "May":
		return time.May
	case "June":
		return time.June
	case "July":
		return time.July
	case "August":
		return time.August
	case "September":
		return time.September
	case "October":
		return time.October
	case "November":
		return time.November
	case "December":
		return time.December
	}
}

func convertQuotes(input string) string {
	input = strings.Replace(input, "’", "'", -1)
	input = strings.Replace(input, "“", "\"", -1)
	input = strings.Replace(input, "”", "\"", -1)
	return input
}

var DA = "([A-Za-z]+)\\s+([A-Za-z]+)\\s+([0-9]+)(st|nd|rd|th),\\s+(20[0-9]{2})"
var dateRegexp = regexp.MustCompile(DA)
var timeRegexp = regexp.MustCompile("([0-9]+):([0-9]+)\\s+(AM|PM)")

func parseDateTime(_date, _time string) time.Time {
	d := dateRegexp.FindStringSubmatch(_date)
	year, _ := strconv.Atoi(d[5])
	day, _ := strconv.Atoi(d[3])
	month := parseMonth(d[2])
	t := timeRegexp.FindStringSubmatch(_time)
	hour, _ := strconv.Atoi(t[1])
	minute, _ := strconv.Atoi(t[2])
	if t[3] == "PM" {
		hour += 12
	}
	dt := time.Date(year, month, day, hour, minute, 0, 0, time.UTC)
	return dt
}

func parseHTML(i int, entry *goquery.Selection) Document {
	document := Document{}

	entryTitle := entry.Find(".entry-title a")
	entryDateTime := entry.Find(".entry-meta a[rel=bookmark]")
	entryTime, _ := entryDateTime.Attr("title")
	entryDate := entryDateTime.Find("time").Text()

	document.URL, _ = entryTitle.Attr("href")
	document.Title = entryTitle.Text()
	document.Date = parseDateTime(entryDate, entryTime)
	document.Image = ""

	body := entry.Next()

	body.Find("style").Each(func(i int, subEntry *goquery.Selection) {
		subEntry.Parent().Get(0).RemoveChild(subEntry.Get(0))
	})

	body.Find("img").Each(func(i int, subEntry *goquery.Selection) {
		image, _ := subEntry.Attr("src")
		if !subEntry.HasClass("wp-smiley") {
			document.Image += image + "\n"
		}
	})

	document.Image = strings.TrimSpace(document.Image)

	content := convertQuotes(strings.TrimSpace(body.Text()))
	if len(content) > 0 {
		document.Content = content
	}

	return document
}

func getTotalPages() int {
	url := fmt.Sprintf(COMMITSTRIP_PAGE_URL, 2)
	debug(url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}
	total, _ := doc.Find(".wp-pagenavi a").Last().Attr("href")

	var totalPages int
	fmt.Sscanf(total, COMMITSTRIP_PAGE_URL, &totalPages)

	return totalPages
}

func request(page int) Document {
	url := COMMITSTRIP_HOME_URL
	if page > 1 {
		url = fmt.Sprintf(COMMITSTRIP_PAGE_URL, page)
	}
	debug(url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}
	document := make(chan Document, 1)
	doc.Find(".entry-header").Each(func(i int, entry *goquery.Selection) {
		document <- parseHTML(i, entry)
	})
	close(document)
	return <-document
}

func batch(documents []Document, start, pages int) {
	promise := make(chan int, pages)
	for i := 0; i < pages; i++ {
		go func(page int) {
			index := start + page
			if index < len(documents) {
				documents[index] = request(start + page + 1)
			}
			promise <- page
		}(i)
	}
	for i := 0; i < pages; i++ {
		<-promise
	}
}

func getDocuments() []Document {
	data, _ := ioutil.ReadFile(DATA_FILE)
	var documents []Document
	json.Unmarshal(data, &documents)
	return documents
}

func addDocuments(oldDocuments *[]Document, newDocuments *[]Document) {
	var isNew bool
	oldDocs := *oldDocuments
	newDocs := *newDocuments

	for i := 0; i < len(newDocs); i++ {
		isNew = true
		for j := 0; j < len(oldDocs); j++ {
			if oldDocs[j].URL == newDocs[i].URL {
				isNew = false
				oldDocs[j].Image = newDocs[i].Image
				if len(oldDocs[j].Content) == 0 && len(newDocs[i].Content) > 0 {
					oldDocs[j].Content = newDocs[i].Content
				}
				break
			}
		}
		if isNew {
			oldDocs = append(oldDocs, newDocs[i])
		}
	}

	sort.Sort(ByDate(oldDocs))
	*oldDocuments = oldDocs
}

func main() {
	// if you want to request for specific page, uncomment code below
	// fmt.Printf("%#v\n", request(355))
	// os.Exit(0)

	var pages int
	pages = getTotalPages()
	debug("Total pages:", pages)

	if len(os.Args) > 1 {
		fmt.Sscanf(os.Args[1], "%d", &pages)
		if pages < 0 {
			pages = 0
		}
	}

	debug("Pages to get:", pages)

	pagesPerBatch := 8
	batches := int(math.Ceil(float64(pages) / float64(pagesPerBatch)))

	newDocuments := make([]Document, pages)

	for i := 0; i < batches; i++ {
		batch(newDocuments, i*pagesPerBatch, pagesPerBatch)
	}

	oldDocuments := getDocuments()

	if oldDocuments == nil || len(oldDocuments) == 0 {
		oldDocuments = newDocuments
	} else {
		addDocuments(&oldDocuments, &newDocuments)
	}

	j, err := json.MarshalIndent(oldDocuments, "", "  ")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(DATA_FILE, append(j, '\n'), 0644)
}
