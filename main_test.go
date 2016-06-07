package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/PuerkitoBio/goquery"
	"log"
	"io"
	"io/ioutil"
	"bytes"
)

func TestExampleCaseViaHttp(t *testing.T) {
	url := "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	results := ExtractResults(doc)
	assert.Equal(t, float32(15.1), results.Total, "The total of that item should be correct")
	assert.Equal(t, 7, len(results.Results), "The length of the results slice should be correct")
	assert.Equal(t, "Apricots", results.Results[0].Description, "The description should be correct")
	assert.Equal(t, "38.3kb", results.Results[0].Size, "The size should be correct")
	assert.Equal(t, "Sainsbury's Apricot Ripe & Ready x5", results.Results[0].Title, "The title should be correct")
	assert.Equal(t, float32(3.5), results.Results[0].UnitPrice, "The unit price should be correct")
}

func TestSingleProductListItem(t *testing.T) {
	reader := loadFile("single_product_list_item.html")
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	results := ExtractResults(doc)

	assert.Equal(t, float32(1.8), results.Total, "The total of that item should be correct")
	assert.Equal(t, 1, len(results.Results), "The length of the results slice should be correct")
	assert.Equal(t, "Gold Kiwi", results.Results[0].Description, "The description should be correct")
	assert.Equal(t, "38.6kb", results.Results[0].Size, "The size should be correct")
	assert.Equal(t, "Sainsbury's Golden Kiwi x4", results.Results[0].Title, "The title should be correct")
	assert.Equal(t, float32(1.8), results.Results[0].UnitPrice, "The unit price should be correct")
}

func TestSingleProductDivProduct(t *testing.T) {
	reader := loadFile("single_product_div_product.html")
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	results := ExtractResults(doc)

	assert.Equal(t, float32(1.8), results.Total, "The total of that item should be correct")
	assert.Equal(t, 1, len(results.Results), "The length of the results slice should be correct")
	assert.Equal(t, "Kiwi", results.Results[0].Description, "The description should be correct")
	assert.Equal(t, "39.0kb", results.Results[0].Size, "The size should be correct")
	assert.Equal(t, "Sainsbury's Kiwi Fruit, Ripe & Ready x4", results.Results[0].Title, "The title should be correct")
	assert.Equal(t, float32(1.8), results.Results[0].UnitPrice, "The unit price should be correct")
}

func TestFullPageOneProduct(t *testing.T) {
	reader := loadFile("full_page_1_product.html")
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	results := ExtractResults(doc)

	assert.Equal(t, float32(3.5), results.Total, "The total of the full page should be correct")
	assert.Equal(t, 1, len(results.Results), "The length of the results slice should be correct")
	assert.Equal(t, "Apricots", results.Results[0].Description, "The description should be correct")
	assert.Equal(t, "38.3kb", results.Results[0].Size, "The size should be correct")
	assert.Equal(t, "Sainsbury's Apricot Ripe & Ready x5", results.Results[0].Title, "The title should be correct")
	assert.Equal(t, float32(3.5), results.Results[0].UnitPrice, "The unit price should be correct")
}

func TestFullPageTenProducts(t *testing.T) {
	reader := loadFile("full_page_10_products.html")
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	results := ExtractResults(doc)

	assert.Equal(t, float32(21.899998), results.Total, "The total of that item should be correct")
	assert.Equal(t, 10, len(results.Results), "The length of the results slice should be correct")
	assert.Equal(t, "Apricots", results.Results[0].Description, "The description should be correct")
	assert.Equal(t, "38.3kb", results.Results[0].Size, "The size should be correct")
	assert.Equal(t, "Sainsbury's Apricot Ripe & Ready x5", results.Results[0].Title, "The title should be correct")
	assert.Equal(t, float32(3.5), results.Results[0].UnitPrice, "The unit price should be correct")
	assert.Equal(t, "Gold Kiwi", results.Results[9].Description, "The description should be correct")
	assert.Equal(t, "38.6kb", results.Results[9].Size, "The size should be correct")
	assert.Equal(t, "Sainsbury's Golden Kiwi x4", results.Results[9].Title, "The title should be correct")
	assert.Equal(t, float32(1.8), results.Results[9].UnitPrice, "The unit price should be correct")
}

func loadFile(name string) io.Reader {
	dat, err := ioutil.ReadFile("test_res/" + name)
	if err != nil {
		panic("Error reading file")
	}
	return bytes.NewReader(dat)
}