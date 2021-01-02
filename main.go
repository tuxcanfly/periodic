package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

type Element struct {
	Name    string
	Image   string
	Content string
}

func main() {
	fName := "elements.json"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}

	elements := make([]Element, 94)

	defer file.Close()
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("www.thoughtco.com"),
	)

	// Extract details of the course
	c.OnHTML(`div.list-sc-item`, func(e *colly.HTMLElement) {
		log.Println("Element found", e.Request.URL)
		title := e.ChildText(`span[class=mntl-sc-block-heading__text]`)
		if title == "" {
			log.Println("No title found", e.Request.URL)
		}
		image := e.ChildAttr(`img`, "data-src")
		if image == "" {
			log.Println("No image found", e.Request.URL)
		}
		content := e.ChildText(`p.mntl-sc-block-html`)
		if image == "" {
			log.Println("No content found", e.Request.URL)
		}
		element := Element{
			Name:    title,
			Image:   image,
			Content: content,
		}
		elements = append(elements, element)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://www.thoughtco.com/chemical-element-pictures-photo-gallery-4052466")

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(elements)
}
