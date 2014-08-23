package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"time"
)

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

var DA = "([A-Za-z]+)\\s([A-Za-z]+)\\s([0-9]+)(st|nd|rd|th),\\s(20[0-9]{2})"
var dateRegexp = regexp.MustCompile(DA)
var timeRegexp = regexp.MustCompile("([0-9]+):([0-9]+)\\s(AM|PM)")

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

func parseDocument(i int, entry *goquery.Selection) {
	entryTitle := entry.Find(".entry-title a")
	permalink, _ := entryTitle.Attr("href")
	title := entryTitle.Text()
	entryDateTime := entry.Find(".entry-meta a[rel=bookmark]")
	entryTime, _ := entryDateTime.Attr("title")
	entryDate := entryDateTime.Find("time").Text()
	dt := parseDateTime(entryDate, entryTime)
	imageSrc, _ := entry.Next().Find("img").Attr("src")
	fmt.Println("Title: ", title)
	fmt.Println("Date:  ", dt)
	fmt.Println("Image: ", imageSrc)
	fmt.Println("URL:   ", permalink)
}

func main() {
	doc, err := goquery.NewDocument("http://www.commitstrip.com/en/")

	if err != nil {
		panic(err)
	}

	doc.Find(".entry-header").Each(parseDocument)
}
