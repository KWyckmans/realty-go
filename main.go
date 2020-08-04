package main

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	const root = "https://www.era.be"
	const startPrice = 200000
	const maxPrice = 400000

	c := colly.NewCollector(
	// colly.AllowedDomains("era.be"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*era.*",
		Parallelism: 4,
		Delay:       1 * time.Second,
	})

	c.OnHTML(".era-search--result-nodes .node-property .field-name-node-link a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// fmt.Println("found:", link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML(".pager li .last", func(e *colly.HTMLElement) {
		// fmt.Println(e.Attr("href"))
		nextPage := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(nextPage))
	})

	c.OnHTML(".intro", func(e *colly.HTMLElement) {
		title := e.ChildText("h1")
		address := e.ChildText(".field-name-era-adres--c a[href]")
		var maps, err = url.Parse(e.ChildAttr(".field-name-era-adres--c a[href]", "href"))

		var bedrooms int
		var bathrooms int
		var livingArea int
		var price int

		rawBathrooms := e.ChildText(".field-name-era-aantal-badkamers--c .era-tooltip-field")
		if len(rawBathrooms) > 0 {
			bathrooms, err = strconv.Atoi(rawBathrooms)
		}

		rawBedrooms := e.ChildText(".field-name-era-aantal-slaapkamers--c .era-tooltip-field")
		if len(rawBedrooms) > 0 {
			bedrooms, err = strconv.Atoi(rawBedrooms)
		}

		rawLivingArea := strings.TrimSpace(strings.Split(e.ChildText(".field-name-era-oppervlakte-bewoonbaar--c .era-tooltip-field"), "m")[0])
		if len(rawLivingArea) > 0 {
			livingArea, err = strconv.Atoi(rawLivingArea)
		}

		rawPrice := e.ChildText(".field-name-era-actuele-vraagprijs--c .field-item")
		if len(rawPrice) > 0 {
			rawPrice = strings.Replace(strings.Split(rawPrice, "€")[1], " ", "", -1)
			price, err = strconv.Atoi(rawPrice)
		} else {
			price = math.MaxInt64
		}

		if err != nil {
			panic(err)
		}

		fmt.Println(title, address, maps, bathrooms, bedrooms, livingArea, price)
	})

	// Find and visit all links

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Println("error:", e, r.Request.URL, string(r.Body))
	})

	var url = fmt.Sprintf("%s/nl/te-koop/wilrijk?broker_id=52636&price=%d+%d", root, startPrice, maxPrice)
	// var url = "https://www.era.be/nl/te-koop/wilrijk/appartement/luxe-appartement-3slpk-dubbele-garagebox-en-staanplaats?broker_id=52636"
	c.Visit(url)
}
