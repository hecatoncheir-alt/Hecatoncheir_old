package mvideo

import (
	"regexp"
	"strconv"
	"time"

	"fmt"

	"github.com/gocolly/colly"
	"github.com/hecatoncheir/Hecatoncheir/crawler"
	"github.com/hecatoncheir/Logger"
)

var cities = crawler.Cities{
	"Москва":      "CityCZ_975",
	"Новосибирск": "CityCZ_2246",
}

// Crawler for parse documents
type Crawler struct {
	LogWriter          logger.Writer
	patternForCutPrice *regexp.Regexp
}

// New create a new Crawler object
func New(logWriter logger.Writer) *Crawler {
	patternForCutPrice, err := regexp.Compile("[ ¤]*")
	if err != nil {
		panic(err)
	}

	parser := Crawler{
		LogWriter:          logWriter,
		patternForCutPrice: patternForCutPrice}
	return &parser
}

func (parser *Crawler) getPagesCount(config crawler.ParserOfCompanyInstructions) (pagesCount int, err error) {
	collector := colly.NewCollector(colly.Async(true))

	collector.OnHTML(config.PageInstruction.PageInPaginationSelector,
		func(element *colly.HTMLElement) {
			pagesCount, err = strconv.Atoi(element.Text)

			if err != nil {
				info := fmt.Sprintf(
					"Get count of pages from: %v failed with response: %v. Error: %v",
					element.Request.URL,
					element.Response.Body,
					err)

				data := logger.LogData{Message: info, Level: "warning"}
				go parser.LogWriter.Write(data)
			}
		})

	collector.OnError(func(response *colly.Response, err error) {
		warning := fmt.Sprintf(
			"Request URL: %v failed with response: %v \nError: %v",
			response.Request.URL,
			response,
			err)

		data := logger.LogData{Message: warning, Level: "warning"}
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

	return pagesCount, nil
}

func (parser *Crawler) getProductsFromPage(
	instructions crawler.ParserOfCompanyInstructions,
	pageIRI string,
	productsChannel chan<- crawler.Product,
	pageIsParsedChannel chan bool) {

	collector := colly.NewCollector()

	collector.OnHTML(instructions.PageInstruction.ItemSelector,
		func(element *colly.HTMLElement) {

			productName := element.ChildText(instructions.PageInstruction.NameOfItemSelector)
			productIRI := instructions.Company.IRI +
				element.ChildAttr(instructions.PageInstruction.LinkOfItemSelector, "href")
			previewImageLink := element.ChildAttr(
				instructions.PageInstruction.PreviewImageOfItemSelector, "data-original")

			product := crawler.Product{
				Language:         instructions.Language,
				Name:             productName,
				IRI:              productIRI,
				City:             instructions.City,
				Category:         instructions.Category,
				Company:          instructions.Company,
				PreviewImageLink: previewImageLink}

			priceOfItemValue := element.ChildText(instructions.PageInstruction.PriceOfItemSelector)
			priceIsMatched := parser.patternForCutPrice.MatchString(priceOfItemValue)

			var priceOfItem string
			if priceIsMatched {
				priceOfItem = parser.patternForCutPrice.ReplaceAllString(priceOfItemValue, "")
			}

			priceValue, err := strconv.ParseFloat(priceOfItem, 64)
			if err != nil {
				warning := fmt.Sprintf(
					"Error get price of product: %v, by IRI: %v",
					product,
					element.Request.URL)

				data := logger.LogData{Message: warning, Level: "warning"}
				go parser.LogWriter.Write(data)
			}

			price := crawler.Price{
				City:     instructions.City,
				Value:    priceValue,
				DateTime: time.Now().UTC()}

			product.Price = price

			info := fmt.Sprintf(
				"Get product for category: %v of compnay: %v by iri: %v",
				instructions.Category,
				instructions.Company,
				element.Request.URL)

			data := logger.LogData{Message: info, Level: "info"}
			go parser.LogWriter.Write(data)
			productsChannel <- product
		})

	collector.OnError(func(response *colly.Response, err error) {
		warning := fmt.Sprintf(
			"Request URL: %v failed with response: %v. Error: %v",
			response.Request.URL,
			response,
			err)

		data := logger.LogData{Message: warning, Level: "warning"}
		go parser.LogWriter.Write(data)
	})

	err := collector.Visit(pageIRI)
	if err != nil {
		warning := fmt.Sprintf(
			"Error visit URL: %v. Error: %v",
			pageIRI,
			err)

		data := logger.LogData{Message: warning, Level: "warning"}
		go parser.LogWriter.Write(data)
	}

	collector.Wait()
	pageIsParsedChannel <- true
}

// RunWithConfiguration can parse web documents and make Item structure
// for each product on page filtered by selectors.
func (parser *Crawler) RunWithConfiguration(config crawler.ParserOfCompanyInstructions) (
	productsChannel chan crawler.Product, err error) {

	productsChannel = make(chan crawler.Product)

	cityCode, err := cities.SearchCodeByCityName(config.City.Name)
	if err != nil {
		return productsChannel, err
	}

	pagesCount, err := parser.getPagesCount(config)
	if err != nil {
		return productsChannel, err
	}

	pageIsParsedChannel := make(chan bool)
	pageConfig := config.PageInstruction
	for pageNumber := 1; pageNumber < pagesCount; pageNumber++ {

		pageIRI := config.Company.IRI +
			pageConfig.Path +
			pageConfig.PageParamPath +
			strconv.Itoa(pageNumber) +
			pageConfig.CityParamPath +
			cityCode

		go parser.getProductsFromPage(config, pageIRI, productsChannel, pageIsParsedChannel)
	}

	parsedPagesCount := 0
	go func() {
		for isPageParsed := range pageIsParsedChannel {
			if isPageParsed {
				parsedPagesCount++
				if parsedPagesCount == pagesCount {
					close(productsChannel)
					break
				}
			} else {
				pageWithError := parsedPagesCount + 1
				warning := fmt.Sprintf(
					"Page is not parsed: %v parsing", pageWithError)

				data := logger.LogData{Message: warning, Level: "warning"}
				go parser.LogWriter.Write(data)
				break
			}
		}
	}()

	return productsChannel, nil
}
