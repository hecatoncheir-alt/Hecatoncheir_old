package mvideo

import (
	"log"
	"regexp"
	"strings"
	"time"

	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/hecatoncheir/Hecatoncheir/crawler"
	"github.com/hecatoncheir/Hecatoncheir/logger"
)

var cities = crawler.Cities{
	"Москва":      "CityCZ_975",
	"Новосибирск": "CityCZ_2246",
}

// Crawler for parse documents
type Crawler struct {
	Items              chan crawler.Product // For subscribe to events
	LogWriter          logger.Writer
	patternForCutPrice *regexp.Regexp
}

// New create a new Crawler object
func New(logWriter logger.Writer) *Crawler {
	patternForCutPrice, err := regexp.Compile("р[уб]*?.")
	if err != nil {
		panic(err)
	}

	parser := Crawler{
		Items:              make(chan crawler.Product),
		LogWriter:          logWriter,
		patternForCutPrice: patternForCutPrice}
	return &parser
}

// getItemsFromPage can get product from html document by selectors in the configuration
func (parser *Crawler) getItemsFromPage(document *goquery.Document, config crawler.ParserOfCompanyInstructions, patternForCutPrice *regexp.Regexp) error {
	pageConfig := config.PageInstruction

	document.Find(pageConfig.ItemSelector).Each(func(iterator int, item *goquery.Selection) {
		var name, price, link, previewImageLink string

		name = item.Find(pageConfig.NameOfItemSelector).Text()
		price = item.Find(pageConfig.PriceOfItemSelector).Text()
		link = item.Find(pageConfig.LinkOfItemSelector).AttrOr("href", "/")
		previewImageLink = item.Find(pageConfig.PreviewImageOfItemSelector).AttrOr("data-original", "/")

		name = strings.TrimSpace(name)
		price = strings.TrimSpace(price)
		link = config.Company.IRI + link
		previewImageLink = strings.Replace(previewImageLink, "//", "", 1)

		// price = strings.Replace(price, "р.", "", -1)
		price = patternForCutPrice.ReplaceAllString(price, "")

		//fmt.Printf("Review %s: %s \n", name, price)

		priceData := crawler.Price{
			Value:    price,
			City:     config.City,
			DateTime: time.Now().UTC(),
		}

		pageItem := crawler.Product{
			Name:             name,
			Price:            priceData,
			IRI:              link,
			Company:          config.Company,
			Language:         config.Language,
			City:             config.City,
			Category:         config.Category,
			PreviewImageLink: previewImageLink,
		}

		log.Println(fmt.Sprintf("Product: '%v' of category: '%v' of company: '%v' parsed. Price: '%s'", name, config.Category.Name, config.Company.Name, priceData.Value))

		parser.Items <- pageItem
	})

	return nil
}

// RunWithConfiguration can parse web documents and make Item structure for each product on page filtered by selectors
func (parser *Crawler) RunWithConfiguration(config crawler.ParserOfCompanyInstructions) error {

	pageConfig := config.PageInstruction

	cityCode, err := cities.SearchCodeByCityName(config.City.Name)
	if err != nil {
		return err
	}

	collector := colly.NewCollector()

	collector.OnRequest(func(request *colly.Request) {
	})

	collector.OnError(func(r *colly.Response, err error) {
		info := fmt.Sprintf("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		data := logger.LogData{Message: info}
		go parser.LogWriter.Write(data)
	})

	firstPageIRI := config.Company.IRI + pageConfig.Path + pageConfig.PageParamPath + "1" + pageConfig.CityParamPath + cityCode

	collector.Visit(firstPageIRI)

	// document, err := goquery.NewDocument(config.Company.IRI + pageConfig.Path + pageConfig.PageParamPath + "1" + pageConfig.CityParamPath + cityCode)
	// if err != nil {
	// 	return err
	// }

	//go parser.GetItemsFromPage(document, config, patternForCutPrice)

	//pagesCount := document.Find(pageConfig.PageInPaginationSelector).Last().Find("a").Text()

	//countOfPages, err := strconv.Atoi(pagesCount)
	//if err != nil {
	//	return err
	//}

	//maxProductsForChannel := 6
	//
	//pagesCrawling := make(chan func(), maxProductsForChannel)
	//
	//go func() {
	//	for craw := range pagesCrawling {
	//		go craw()
	//	}
	//}()
	//
	//// Because first page already parsed for get pages count
	//pageNumberFromCrawlingStart := 2
	//
	//var iterator int
	//for iterator = pageNumberFromCrawlingStart; iterator <= countOfPages; iterator++ {
	//	document, err := goquery.NewDocument(config.Company.IRI + pageConfig.Path + pageConfig.PageParamPath + strconv.Itoa(iterator) + pageConfig.CityParamPath + cityCode)
	//	if err != nil {
	//		return err
	//	}
	//
	//	pagesCrawling <- func() {
	//		parser.GetItemsFromPage(document, config, patternForCutPrice)
	//	}
	//}
	//
	//close(pagesCrawling)

	return nil
}
