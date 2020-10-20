package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	const startPrice = 200000
	const maxPrice = 400000

	log.Println("Load properties")
	var properties []Property = loadProperties("test.json")

	log.Println("Start scraping")
	properties = scrapeEra(properties, startPrice, maxPrice)

	log.Println("Finished scraping, saving to file")

	// This will dump everything in proprties to a file. Saving everything everytime
	// seems awefully inefficient, but it will have to do for now.
	// An easy optimisation would be to check if I even found any updates, but that will
	// happen very often, so only has a minor impact.
	saveResults(properties, "test.json")
}

/**
Checks if we already found any given property based on the url where we found it.
This works for now, but will not work once a second website for scraping is added.

This is also the reason I'm not using maps to store the properties in for now. I could
store properties in a map, using the url as key, but that again would break down once a second
site is added that may contain the same properties (and I know for a fact that that will happen).

I could still use this simple find within results for one given site. Potential optimisation.
*/
func findProperty(properties []Property, url url.URL) int {
	for i, p := range properties {
		if url == p.URL {
			return i
		}
	}

	return -1
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
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnHTML(".pager li .last", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(nextPage))
	})

	c.OnHTML(".intro", func(e *colly.HTMLElement) {
		if findProperty(properties, *e.Request.URL) == -1 {
			property := parseProperty(e)
			properties = append(properties, property)
		} else {
			log.Println("Encountered property that was already known")
		}
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Println("error:", e, r.Request.URL, string(r.Body))
	})

	var url = fmt.Sprintf("%s/nl/te-koop/wilrijk?broker_id=52636&price=%d+%d", root, startPrice, maxPrice)

	c.Visit(url)

	return properties
}

func loadProperties(filename string) []Property {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var properties []Property
	err = json.Unmarshal(data, &properties)

	return properties
}

func saveResults(properties []Property, filename string) {
	file, _ := json.MarshalIndent(properties, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)
}
