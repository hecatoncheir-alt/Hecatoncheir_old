package main

import (
	"github.com/hecatoncheir/Hecatoncheir/broker"
	"github.com/hecatoncheir/Hecatoncheir/configuration"

	"encoding/json"
	"testing"

	"github.com/hecatoncheir/Hecatoncheir/crawler"
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

	parserOfCompany := crawler.ParserOfCompanyInstructions{
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

		go handlesNeedProductsOfCategoryOfCompanyEvent(data["Data"], bro, "test")

		break
	}

	channelForGetProducts, err := bro.ListenTopic("test", "Parser")
	if err != nil {
		test.Error(err)
	}

	for message := range channelForGetProducts {
		data := map[string]interface{}{}
		json.Unmarshal(message, &data)


		if data["Message"] != "Product of category of company ready" {
			test.Fail()
		}

		if data["Data"].(map[string]interface{})["Language"] != "en" {
			test.Fail()
		}

		break
	}
}
