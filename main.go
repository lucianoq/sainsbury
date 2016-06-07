package main

import (
	"net/http"
	"log"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"encoding/json"
)

const URL = "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"

type Output struct {
	Results []Item `json:"results"`
	Total   float32 `json:"total"`
}

type Item struct {
	Title       string `json:"title"`
	Size        string `json:"size"`
	UnitPrice   float32 `json:"unit_price"`
	Description string `json:"description"`
}

func main() {
	// Download and parse the document
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		log.Fatal(err)
	}

	results := ExtractResults(doc)

	// Print the output
	jsonString, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", jsonString)
}


// Build the results
func ExtractResults(doc *goquery.Document) Output {
	var output Output

	// for each product in the page, extract info and add it to the results slice
	doc.Find(".product").Each(func(i int, s *goquery.Selection) {
		response := getHttpResponseFromProduct(s)
		item := Item{
			Title: ExtractTitle(s),
			UnitPrice: ExtractUnitPrice(s),
			Description: ExtractDescription(response),
			Size: ExtractSize(response),
		}
		output.Results = append(output.Results, item)
		output.Total += item.UnitPrice
	})

	return output
}


// Retrieve title
func ExtractTitle(product *goquery.Selection) string {
	return strings.TrimSpace(product.Find(".productInfo").Find("a").Text())
}


// Retrieve the unit_price manipulating and parsing the string
func ExtractUnitPrice(product *goquery.Selection) float32 {
	priceString := strings.TrimSpace(product.Find(".pricePerUnit").Text())
	priceString = strings.Trim(priceString, "/unit£")
	unitPrice, err := strconv.ParseFloat(priceString, 32)
	if err != nil {
		log.Fatal(err)
	}
	return float32(unitPrice)
}


// Retrieve the size of the page
// from the contentLength of the http response
func ExtractSize(response *http.Response) string {
	return fmt.Sprintf("%.1fkb", float32(response.ContentLength) / 1024)
}


// Retrieve description
func ExtractDescription(response *http.Response) string {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		log.Fatal(err)
	}

	return doc.Find("#information").Find(".productDataItemHeader").FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Text() == "Description"
	}).Next().Children().First().Text()
}


// Retrieve the http Response from the page linked by title
// using http.Get() instead of goquery.NewDocument() to avoid a double http call
// to have both size and content
func getHttpResponseFromProduct(product *goquery.Selection) *http.Response {
	link := product.Find(".productInfo").Find("a")
	url, exists := link.Attr("href")
	if !exists {
		log.Fatal("href does not exists in a tag")
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}