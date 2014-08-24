package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math"
	"regexp"
	"strconv"
	"time"
)

type Document struct {
	URL   string
	Title string
	Date  time.Time
	Image string
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
	document.Image, _ = entry.Next().Find("img").Attr("src")

	return document
}

func request(page int) Document {
	url := "http://www.commitstrip.com/en/"
	if page > 1 {
		url = fmt.Sprintf("%spage/%d/", url, page)
	}
	fmt.Println("DEBUG", url)
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

func main() {
	pages := 9
	pagesPerBatch := 8
	batches := int(math.Ceil(float64(pages) / float64(pagesPerBatch)))

	documents := make([]Document, pages)

	for i := 0; i < batches; i++ {
		batch(documents, i*pagesPerBatch, pagesPerBatch)
	}

	for i, document := range documents {
		if i > 0 {
			fmt.Println()
		}
		fmt.Println("Title:", document.Title)
		fmt.Println("Date: ", document.Date)
		fmt.Println("Image:", document.Image)
		fmt.Println("URL:  ", document.URL)
	}
}
