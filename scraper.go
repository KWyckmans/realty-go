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

func scrapeEra(properties []Property, startPrice int, maxPrice int) []Property {
	const root = "https://www.era.be"

	c := colly.NewCollector()

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*era.*",
		Parallelism: 4,
		Delay:       1 * time.Second,
	})

	c.OnHTML(".era-search--result-nodes .node-property .field-name-node-link a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// log.Println("visiting ", link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML(".pager li .last", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
		// log.Println("Going to next page")
		c.Visit(e.Request.AbsoluteURL(nextPage))
	})

	c.OnHTML(".intro", func(e *colly.HTMLElement) {
		i := findProperty(properties, *e.Request.URL)
		if i == -1 {
			property := parseProperty(e)
			properties = append(properties, property)
			i = len(properties) - 1
		} else {
			properties[i].LastUpdated = time.Now()
		}

		// We check both sold an options on properties
		// on both old and new properties.
		// This to cover any holes in scraping:
		// Say we haven't scraped in a couple of days, in that case
		// properties may appear that are "directly" sold.
		// Same in the other direction: we may have old properties
		// stored taht have an option next time we scrape.
		isSold := verifyIfSold(e)
		isOption := verifyIfOption(e)

		if isSold && properties[i].SoldAt.IsZero() {
			properties[i].SoldAt = time.Now()
		}

		if isOption && properties[i].OptionAt.IsZero() {
			properties[i].OptionAt = time.Now()
		}
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Println("error:", e, r.Request.URL, string(r.Body))
	})

	var url = fmt.Sprintf("%s/nl/te-koop/wilrijk?broker_id=52636&price=%d+%d", root, startPrice, maxPrice)

	c.Visit(url)

	return properties
}

func parseProperty(e *colly.HTMLElement) Property {
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
	if len(rawPrice) > 0 && rawPrice != "Op aanvraag" {
		rawPrice = strings.Replace(strings.Split(rawPrice, "€")[1], " ", "", -1)
		price, err = strconv.Atoi(rawPrice)
	} else {
		price = math.MaxInt64
	}

	if err != nil {
		log.Println("an error occured!", err)
	}

	return NewProperty(price, livingArea, bedrooms, bathrooms, address, *maps, title, *e.Request.URL)
}

/**
We can potentially miss properties that are sold if the program hasn't run everyday:
We suddenly see a proprty that starts at the sold state. Assuming they won't put up
properties that are already sold
*/
func verifyIfSold(e *colly.HTMLElement) bool {
	rawSold := e.ChildText(".property-state span")
	if len(rawSold) > 0 && rawSold == "Verkocht" {
		return true
	}

	return false
}

/**
We can potentially miss properties that are sold if the program hasn't run everyday:
We suddenly see a proprty that starts at the sold state. Assuming they won't put up
properties that are already sold
*/
func verifyIfOption(e *colly.HTMLElement) bool {
	rawSold := e.ChildText(".property-state span")
	if len(rawSold) > 0 && rawSold == "In optie" {
		return true
	}

	return false
}
