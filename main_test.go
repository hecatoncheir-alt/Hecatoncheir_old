package main

import (
	"github.com/hecatoncheir/Hecatoncheir/broker"
	"github.com/hecatoncheir/Hecatoncheir/configuration"

	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Hecatoncheir/crawler"
	"testing"
)

func TestIntegrationCanParseCategoryOfCompanyByBrokerEventRequest(test *testing.T) {
	config, err := configuration.GetConfiguration()
	if err != nil {
		test.Error(err)
	}

	bro := broker.New()

	err = bro.Connect(config.Development.Broker.Host, config.Development.Broker.Port)
	if err != nil {
		test.Error(err)
	}

	channel, err := bro.ListenTopic("test", "Parser")
	if err != nil {
		test.Error(err)
	}

	parserOfCompany := crawler.ParserOfCompany{
		Language: "en",
		Company: crawler.Company{
			ID:   "0x2786",
			Name: "M.Video",
			IRI:  "http://www.mvideo.ru/"},
		Category: crawler.Category{
			ID:   "",
			Name: "Test category of M.Video company"},
		City: crawler.City{
			ID:   "0x2788",
			Name: "Москва"},
		PageInstruction: crawler.PageInstruction{
			ID:   "0x2789",
			Path: "smartfony-i-svyaz/smartfony-205",
			PageInPaginationSelector:   ".pagination-list .pagination-item",
			PreviewImageOfItemSelector: ".product-tile-picture-link img",
			PageParamPath:              "/f/page=",
			CityParamPath:              "?cityId=",
			CityParam:                  "CityCZ_975",
			ItemSelector:               ".grid-view .product-tile",
			NameOfItemSelector:         ".product-tile-title",
			LinkOfItemSelector:         ".product-tile-title a",
			PriceOfItemSelector:        ".product-price-current"},
	}

	parseData, err := json.Marshal(parserOfCompany)
	if err != nil {
		test.Error(err)
	}

	go bro.WriteToTopic("test", map[string]interface{}{"Message": "Need products of category of company", "Data": string(parseData)})

	for message := range channel {
		data := map[string]string{}
		json.Unmarshal(message, &data)

		if data["Message"] != "Need products of category of company" {
			test.Fail()
		}

		parserOfCompany, err := crawler.NewParserFromJSON(data["Data"])
		if err != nil {
			test.Error(err)
		}

		channelWithProducts, err := parserOfCompany.ReadProductsFromCategoryOfCompany()
		if err != nil {
			test.Error(err)
		}

		for product := range channelWithProducts {
			//TODO
			parsedProduct, err := json.Marshal(product)
			if err != nil {
				test.Error(err)
			}

			go bro.WriteToTopic("test", map[string]interface{}{"Message": "Product of category of company ready", "Data": string(parsedProduct)})
			break
		}

		break
	}

	for message := range channel {
		data := map[string]interface{}{}
		json.Unmarshal(message, &data)

		fmt.Println(data["Message"])

		break
	}

}
