package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
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
	saveProperties(properties, "test.json")
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
