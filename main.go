package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

func main() {
	var isQueryMode = flag.Bool("query", false, "Inidcate wether you want to query the program or just have it run")
	flag.Parse()

	f, err := os.OpenFile("realty.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	log.Println("Loading properties")
	var properties []Property = loadProperties("test.json")
	log.Println(len(properties), "properties loaded")

	if *isQueryMode {
		log.Println("Starting query mode")

		var i int
		for {
			log.Println("Please select an option: ")
			log.Println("1. Print properties")
			log.Println("9. Quit")
			fmt.Scan(&i)

			switch i {
			case 1:
				log.Println("Available properties:")
			case 2:
				log.Println("Showing cheapest properties:")
			}

			if i == 9 {
				break
			}
		}
	} else {
		const startPrice = 200000
		const maxPrice = 400000

		log.Println("Start scraping")
		properties = scrapeEra(properties, startPrice, maxPrice)

		log.Println("Finished scraping, saving to file")

		// This will dump everything in proprties to a file. Saving everything everytime
		// seems awefully inefficient, but it will have to do for now.
		// An easy optimisation would be to check if I even found any updates, but that will
		// happen very often, so only has a minor impact.
		saveProperties(properties, "test.json")
	}
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

func loadProperties(filename string) []Property {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var properties []Property
	err = json.Unmarshal(data, &properties)

	return properties
}

func saveProperties(properties []Property, filename string) {
	file, _ := json.MarshalIndent(properties, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)
}
