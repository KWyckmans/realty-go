package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	ROOT := "https://www.era.be"
	START_PRICE := 200000
	MAX_PRICE := 400000

	c := colly.NewCollector(
	// colly.AllowedDomains("era.be"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*era.*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	c.OnHTML(".node-property .field-name-node-link a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println(link)
		// c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML(".pager li .last", func(e *colly.HTMLElement) {
		fmt.Println(e.Attr("href"))
		nextPage := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(nextPage))
	})

	// Find and visit all links

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Println("error:", e, r.Request.URL, string(r.Body))
	})

	var url = fmt.Sprintf("%s/nl/te-koop/wilrijk?broker_id=52636&price=%d+%d", ROOT, START_PRICE, MAX_PRICE)
	c.Visit(url)
}
