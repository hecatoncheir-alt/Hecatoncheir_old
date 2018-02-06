package mvideo

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/hecatoncheir/Hecatoncheir/crawler"
)

var cities = crawler.Cities{
	"Москва":      "CityCZ_975",
	"Новосибирск": "CityCZ_2246",
}

// Crawler for parse documents
type Crawler struct {
	Items chan crawler.Product // For subscribe to events
}

// NewCrawler create a new Crawler object
func NewCrawler() *Crawler {
	parser := Crawler{Items: make(chan crawler.Product)}
	return &parser
}

// GetItemsFromPage can get product from html document by selectors in the configuration
func (parser *Crawler) GetItemsFromPage(document *goquery.Document, pageConfig crawler.PageInstruction, company crawler.Company, patternForCutPrice *regexp.Regexp) error {
	document.Find(pageConfig.ItemSelector).Each(func(iterator int, item *goquery.Selection) {
		var name, price, link, previewImageLink string

		name = item.Find(pageConfig.NameOfItemSelector).Text()
		price = item.Find(pageConfig.PriceOfItemSelector).Text()
		link = item.Find(pageConfig.LinkOfItemSelector).AttrOr("href", "/")
		previewImageLink = item.Find(pageConfig.PreviewImageOfItemSelector).AttrOr("data-original", "/")

		name = strings.TrimSpace(name)
		price = strings.TrimSpace(price)
		link = company.IRI + link
		previewImageLink = strings.Replace(previewImageLink, "//", "", 1)

		// price = strings.Replace(price, "р.", "", -1)
		price = patternForCutPrice.ReplaceAllString(price, "")

		//fmt.Printf("Review %s: %s \n", name, price)

		cityName, err := cities.SearchCityByCode(pageConfig.CityParam)
		if err != nil {
			log.Println(err)
		}

		priceData := crawler.Price{
			Value:    price,
			City:     cityName,
			DateTime: time.Now().UTC(),
		}

		pageItem := crawler.Product{
			Name:             name,
			Price:            priceData,
			IRI:              link,
			Company:          company,
			PreviewImageLink: previewImageLink,
		}

		parser.Items <- pageItem
	})

	return nil
}

// RunWithConfiguration can parse web documents and make Item structure for each product on page filtered by selectors
func (parser *Crawler) RunWithConfiguration(config crawler.ParserOfCompany) error {
	patternForCutPrice, _ := regexp.Compile("р[уб]*?.")

	pageConfig := config.PageInstruction

	document, err := goquery.NewDocument(config.Company.IRI + pageConfig.Path + pageConfig.PageParamPath + "1" + pageConfig.CityParamPath + pageConfig.CityParam)
	if err != nil {
		return err
	}

	go parser.GetItemsFromPage(document, pageConfig, config.Company, patternForCutPrice)

	pagesCount := document.Find(pageConfig.PageInPaginationSelector).Last().Find("a").Text()

	countOfPages, err := strconv.Atoi(pagesCount)
	if err != nil {
		return err
	}

	maxProductsForChannel := 6

	pagesCrawling := make(chan func(), maxProductsForChannel)

	go func() {
		for craw := range pagesCrawling {
			go craw()
		}
	}()

	// Because first page already parsed for get pages count
	pageNumberFromCrawlingStart := 2

	var iterator int
	for iterator = pageNumberFromCrawlingStart; iterator <= countOfPages; iterator++ {
		document, err := goquery.NewDocument(config.Company.IRI + pageConfig.Path + pageConfig.PageParamPath + strconv.Itoa(iterator) + pageConfig.CityParamPath + pageConfig.CityParam)
		if err != nil {
			return err
		}

		pagesCrawling <- func() {
			parser.GetItemsFromPage(document, pageConfig, config.Company, patternForCutPrice)
		}
	}

	close(pagesCrawling)

	return nil
}
