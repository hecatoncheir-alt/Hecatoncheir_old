package mvideo

import (
	"errors"
	"log"
	"regexp"
	"strconv"
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

var ErrGetPagesCountFail = errors.New("get pages count fail")

func (parser *Crawler) getPagesCount(config crawler.ParserOfCompanyInstructions) (pagesCount int, err error) {
	collector := colly.NewCollector(colly.Async(true))

	collector.OnHTML(config.PageInstruction.PageInPaginationSelector,
		func(element *colly.HTMLElement) {
			pagesCount, err = strconv.Atoi(element.Text)

			if err != nil {
				info := fmt.Sprintf(
					"Get count of pages from: %v failed with response: %v \nError: %v",
					element.Request.URL,
					element.Response.Body,
					err)

				data := logger.LogData{Message: info, Level: "warning"}
				go parser.LogWriter.Write(data)
			}
		})

	collector.OnError(func(response *colly.Response, err error) {
		info := fmt.Sprintf(
			"Request URL: %v failed with response: %v \nError: %v",
			response.Request.URL,
			response,
			err)

		data := logger.LogData{Message: info, Level: "warning"}
		go parser.LogWriter.Write(data)
	})

	cityCode, err := cities.SearchCodeByCityName(config.City.Name)
	if err != nil {
		return 0, err
	}

	pageConfig := config.PageInstruction
	firstPageIRI := config.Company.IRI + pageConfig.Path + pageConfig.PageParamPath + "1" + pageConfig.CityParamPath + cityCode

	err = collector.Visit(firstPageIRI)
	if err != nil {
		return 0, err
	}

	collector.Wait()

	return pagesCount, ErrGetPagesCountFail
}

// RunWithConfiguration can parse web documents and make Item structure for each product on page filtered by selectors
func (parser *Crawler) RunWithConfiguration(config crawler.ParserOfCompanyInstructions) error {

	pagesCount, err := parser.getPagesCount(config)
	if err != nil {
		return err
	}

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
