package main

import (
	"net/url"
)

type address struct {
	address string
	maps    url.URL
}

// Property represents a house or an apparetment found on an immo website
type Property struct {
	Bedrooms   int
	Bathrooms  int
	LivingArea int
	Price      int
	Address    address
	Name       string
	URL        url.URL
}

func (p Property) pricePerSqm() float64 {
	return float64(p.Price) / float64(p.LivingArea)
}

func newAddress(loc string, maps url.URL) address {
	return address{
		address: loc,
		maps:    maps,
	}
}

// NewProperty creates a new property
func NewProperty(price int, livingArea int, bedrooms int, bathrooms int, address string, maps url.URL, name string, url url.URL) Property {
	return Property{
		Price:      price,
		LivingArea: livingArea,
		Bedrooms:   bedrooms,
		Bathrooms:  bathrooms,
		Address:    newAddress(address, maps),
		Name:       name,
		URL:        url,
	}
}
